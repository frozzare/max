package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/frozzare/max/internal/backend"
	"github.com/frozzare/max/internal/backend/docker"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Runner represents a the runner.
type Runner struct {
	args    map[string]interface{}
	ctx     context.Context
	engine  backend.Engine
	Config  *config.Config
	Once    bool
	opts    []Option
	Verbose bool
}

// New creates a new runner.
func New(opts ...Option) *Runner {
	r := &Runner{
		opts: opts,
		ctx:  context.Background(),
	}

	for _, opts := range opts {
		opts(r)
	}

	return r
}

func (r *Runner) exec(t *task.Task) error {
	t = r.prepareTask(t)

	// Use docker if docker configuration is not nil.
	if t.Docker != nil {
		engine, err := docker.New()
		if err != nil {
			return err
		}
		r.engine = engine
	}

	// Run deps before task.
	for _, id := range t.Deps {
		if err := New(append(r.opts, Once(true))...).Run(id); err != nil {
			return err
		}
	}

	// Run other tasks.
	for _, id := range t.Tasks.Values {
		if err := New(r.opts...).Run(id); err != nil {
			return err
		}
	}

	// Execute task in engine.
	if err := r.engine.Exec(t); err != nil {
		return err
	}

	for {
		exited, err := r.engine.Wait(t)
		if err != nil {
			return err
		}

		if exited {
			break
		}
	}

	// Get logs from engine.
	rc, err := r.engine.Logs(t)
	if err != nil {
		return err
	}

	go func() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(rc)
		log.Print(buf.String())
		rc.Close()
	}()

	return nil
}

func (r *Runner) execAll(t *task.Task) <-chan error {
	var g errgroup.Group
	done := make(chan error)

	g.Go(func() error {
		return r.exec(t)
	})

	go func() {
		done <- g.Wait()
		close(done)
	}()

	return done
}

func (r *Runner) prepareTask(t *task.Task) *task.Task {
	// Merge global variables with task variables.
	for k, v := range r.Config.Variables {
		t.Variables[k] = v
	}

	t.Verbose = r.Verbose

	r.args = r.Config.Args
	r.parseArgs()

	if t.Args == nil {
		t.Args = make(map[string]interface{})
	}

	if len(r.args) > 0 {
		for k, v := range r.args {
			t.Args[k] = v
		}
	}

	return t
}

func (r *Runner) parseArgs() {
	var input string

	if len(os.Args) > 1 {
		input = strings.Join(os.Args[2:], " ")
	} else {
		input = strings.Join(os.Args, " ")
	}

	buff := bytes.NewBufferString(input)

	if r.args == nil {
		r.args = make(map[string]interface{})
	}

	for {
		rn, _, err := buff.ReadRune()

		if err != nil {
			break
		}

		if buff.Len() == 0 {
			break
		}

		if rn == '-' {
			if arg, err := buff.ReadString(' '); err == nil {
				if val, err := buff.ReadString(' '); err == nil || err == io.EOF {
					if arg[0] == '-' {
						arg = arg[1:]
					}

					key := strings.TrimSpace(arg)
					key = strings.Replace(key, "-", "_", -1)

					r.args[key] = val
				}
			}
		}
	}
}

// Task returns a task by name if it exists.
func (r *Runner) Task(name string) *task.Task {
	if r.Config == nil || r.Config.Tasks[name] == nil {
		return nil
	}

	for _, n := range []string{fmt.Sprintf("%s_%s", name, runtime.GOOS), name} {
		if t := r.Config.Tasks[n]; t != nil {
			return t
		}
	}

	return nil
}

// Run runs a task.
func (r *Runner) Run(id string) error {
	task := r.Task(id)

	if task == nil {
		return fmt.Errorf("task missing: %s", id)
	}

	done := make(chan bool, 1)

	task.SetID(id)

	defer func() {
		r.engine.Destroy(task)
	}()

	if err := r.engine.Setup(task); err != nil {
		return err
	}

	var e error

	select {
	case <-r.ctx.Done():
		close(done)
		return errors.New("context is done")
	case err := <-r.execAll(task):
		if err != nil {
			e = err
		}
		done <- true
	}

	// Wait for task to be done.
	for {
		if len(done) == 1 {
			close(done)
			break
		}

		time.Sleep(1 * time.Second)
	}

	return e
}
