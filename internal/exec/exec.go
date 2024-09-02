package exec

import (
	"context"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Options represents execute options.
type Options struct {
	Context context.Context
	Dir     string
	Env     []string
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

// Exec will execute a input cmd string.
func Exec(opts *Options) error {

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}

	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}

	if len(opts.Dir) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		opts.Dir = wd
	}

	p, err := syntax.NewParser().Parse(strings.NewReader(opts.Command), "")
	if err != nil {
		return err
	}

	env := os.Environ()
	env = append(env, opts.Env...)
	r, err := interp.New(
		interp.Env(expand.ListEnviron(env...)),
		interp.StdIO(opts.Stdin, opts.Stdout, opts.Stdout),
		interp.Dir(opts.Dir),
	)

	if err != nil {
		return err
	}

	return r.Run(opts.Context, p)
}
