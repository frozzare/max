package exec

import (
	"os"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Exec will execute a input cmd string.
func Exec(input string, args ...string) error {
	var path string

	if len(args) > 0 {
		path = args[0]
	}

	if len(path) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		path = wd
	}

	p, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		return err
	}

	r := interp.Runner{
		Dir:    path,
		Exec:   interp.DefaultExec,
		Open:   interp.OpenDevImpls(interp.DefaultOpen),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err = r.Reset(); err != nil {
		return err
	}

	return r.Run(p)
}
