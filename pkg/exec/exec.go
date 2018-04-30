package exec

import (
	"os"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Options represents execute options.
type Options struct {
	Dir     string
	Env     []string
	Command string
}

// Exec will execute a input cmd string.
func Exec(opts *Options) error {
	var path string

	if len(opts.Dir) > 0 {
		path = opts.Dir
	}

	if len(path) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		path = wd
	}

	p, err := syntax.NewParser().Parse(strings.NewReader(opts.Command), "")
	if err != nil {
		return err
	}

	env := os.Environ()
	for _, e := range opts.Env {
		env = append(env, e)
	}

	r := interp.Runner{
		Env:    env,
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
