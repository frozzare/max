package cmd

import (
	"github.com/frozzare/go/env"
	"github.com/spf13/pflag"
)

func defaultTask() string {
	if v := env.Get("MAX_DEFAULT_TASK"); len(v) > 0 {
		return v
	}

	return "default"
}

func taskWithArgs() (string, []string) {
	args := pflag.Args()
	if len(args) == 0 {
		return defaultTask(), []string{}
	}

	task := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	return task, args
}
