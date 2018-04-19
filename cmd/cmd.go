package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/runner"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const usage = `
Runs the specified task(s).

Commands:

  help [task]           show task help.
  version               print Max version

Options:

`

func readConfig(path string) (*config.Config, error) {
	var c *config.Config
	var err error

	fi, err := os.Stdin.Stat()
	if fi.Mode()&os.ModeNamedPipe != 0 {
		buf, err := ioutil.ReadAll(os.Stdin)

		if err == nil {
			c, err = config.ReadContent(string(buf))

			if err != nil {
				return nil, errors.Wrap(err, "max")
			}
		}
	} else {
		c, err = config.ReadFile(path)
	}

	if err != nil {
		return nil, errors.Wrap(err, "max")
	}

	if c == nil {
		return nil, errors.New("max: bad config")
	}

	return c, nil
}

// Execute executes the command line.
func Execute(version string) {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		log.Print(usage)
		pflag.PrintDefaults()
	}

	var (
		allFlag     bool
		configFile  string
		listFlag    bool
		onceFlag    bool
		verboseFlag bool
	)

	pflag.BoolVarP(&allFlag, "all", "a", false, "runs all tasks")
	pflag.StringVarP(&configFile, "config", "c", "", "sets the config file")
	pflag.BoolVarP(&listFlag, "list", "l", false, "lists tasks with summary description")
	pflag.BoolVarP(&onceFlag, "once", "o", false, "runs tasks once and ignore interval")
	pflag.BoolVarP(&verboseFlag, "verbose", "v", false, "verbose mode")

	pflag.CommandLine.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags: true,
	}
	pflag.Parse()

	// Find arguments to run.
	args := pflag.Args()
	task := ""
	if len(args) == 0 {
		task = "default"
		args = []string{}
	} else {
		task = args[0]
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = []string{}
		}
	}

	// Output max verison.
	if task == "version" {
		log.Printf("Max version: %s\n", version)
		return
	}

	if task == "help" && len(args) == 0 {
		pflag.PrintDefaults()
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
		fmt.Println(c.Tasks)
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
		}

		if len(t.Usage) != 0 {
			log.Printf("Usage:\n  max %s %s\n", id, t.Usage)
		}

		if len(t.Summary) != 0 {
			log.Printf("Summary:\n  %s", t.Summary)
		}

		return
	}

	runner.Run(task)
}
