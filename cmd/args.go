package cmd

import (
	"github.com/spf13/pflag"
)

func getTaskWithArgs() (string, []string) {
	args := pflag.Args()
	if len(args) == 0 {
		return "default", []string{}
	}

	task := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	return task, args
}
