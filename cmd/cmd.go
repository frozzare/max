package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/frozzare/max/internal/runner"
	"github.com/spf13/pflag"
)

const usage = `
Runs the specified task(s).

Commands:

  cache flush           flush cache.
  help [task]           show task help.
  version               print max version.

Options:

`

// Execute executes the command line.
func Execute() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		log.Print(usage)
		pflag.PrintDefaults()
		log.Println("\nUse \"max help [task]\" for more information about that task.")
	}

	var (
		allFlag     bool
		configFile  string
		listFlag    bool
		onceFlag    bool
		verboseFlag bool
	)

	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	pflag.BoolVarP(&allFlag, "all", "a", false, "runs all tasks")
	pflag.StringVarP(&configFile, "config", "c", "", "sets the config file")
	pflag.BoolVarP(&listFlag, "list", "l", false, "lists tasks with summary description")
	pflag.BoolVarP(&onceFlag, "once", "o", false, "runs tasks once and ignore interval")
	pflag.BoolVarP(&verboseFlag, "verbose", "v", false, "verbose mode")

	pflag.CommandLine.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags: true,
	}
	pflag.Parse()

	// Bail if help flag.
	if os.Args[len(os.Args)-1] == "--help" {
		return
	}

	// Find task and arguments to run.
	task, args := getTaskWithArgs()

	// Run built in commands.
	if runCommands(task, args) {
		return
	}

	// Try to read max config file.
	c, err := readConfig(configFile)
	if err != nil {
		log.Printf("max: %s\n", err.Error())
		return
	}

	// Output list of tasks.
	if listFlag {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		for k, t := range c.Tasks {
			fmt.Fprintf(w, "* %s: \t%s\n", k, t.Summary)
		}
		w.Flush()
		return
	}

	runner := runner.Runner{
		All:     allFlag,
		Config:  c,
		Once:    onceFlag,
		Verbose: verboseFlag,
	}

	// Output help usage if requested.
	if task == "help" && len(args) == 1 {
		id := args[0]
		t := runner.Task(id)

		if t == nil {
			log.Fatalf("Task missing: %s", id)
			return
		}

		t.PrintUsage(id)
		return
	}

	runner.Run(task)
}
