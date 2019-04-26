package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/frozzare/max/internal/config"
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

	var (
		c           *config.Config
		configFile  string
		err         error
		onceFlag    bool
		quietFlag   bool
		verboseFlag bool
	)

	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	pflag.StringVarP(&configFile, "config", "c", "", "sets the config file")
	pflag.BoolVarP(&onceFlag, "once", "o", false, "runs tasks once and ignore interval")
	pflag.BoolVarP(&quietFlag, "quiet", "q", false, "minimal logs")
	pflag.BoolVarP(&verboseFlag, "verbose", "v", false, "verbose logs")
	pflag.Parse()

	pflag.CommandLine.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags: true,
	}

	// Read config file if it exists.
	c, err = readConfig(configFile)
	if err != nil {
		log.Printf("max: %s\n", err.Error())
		return
	}

	pflag.Usage = func() {
		log.Print(usage)
		pflag.PrintDefaults()

		if c != nil {
			log.Printf("\n%s\n\n", "Tasks:")
			l := 21
			w := tabwriter.NewWriter(os.Stdout, 8, 01, 0, '\t', 0)
			for k, t := range c.Tasks {
				s := ""
				for i := 0; i < l-len(k); i++ {
					s += " "
				}
				fmt.Fprintf(w, "  %s:%s%s\n", k, s, t.Summary)
			}
			w.Flush()
		}

		log.Println("\nUse \"max help [task]\" for more information about that task.")
	}

	// Bail if help flag.
	if strings.Contains(os.Args[len(os.Args)-1], "-help") {
		return
	}

	// Find task and arguments to run.
	task, args := taskWithArgs()

	// Run built in commands.
	if runCommands(task, args) {
		return
	}

	// Try to read max config file if nil.
	if c == nil {
		c, err = readConfig(configFile)
		if err != nil {
			log.Printf("max: %s\n", err.Error())
			return
		}
	}

	// Create a new runner.
	runner := runner.New(
		runner.Config(c),
		runner.Once(onceFlag),
		runner.Quiet(quietFlag),
		runner.Verbose(verboseFlag),
	)

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

	// Run and log error.
	if err := runner.Run(task); err != nil {
		log.Fatalf("max: %s", err.Error())
	}
}
