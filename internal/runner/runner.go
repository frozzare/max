package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/frozzare/max/internal/backend"
	backendConfig "github.com/frozzare/max/internal/backend/config"
	"github.com/frozzare/max/internal/backend/docker"
	"github.com/frozzare/max/internal/backend/local"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Runner represents a the runner.
type Runner struct {
	args    map[string]interface{}
	ctx     context.Context
	engine  backend.Engine
	config  *config.Config
	log     *log.Logger
	once    bool
	opts    []Option
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	verbose bool
}

// New creates a new runner.
func New(opts ...Option) *Runner {
	r := &Runner{
		opts:   opts,
		ctx:    context.Background(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	for _, opts := range opts {
		opts(r)
	}

	if r.log == nil {
		r.log = log.New(os.Stderr, "", log.LstdFlags)
	}

	return r
}

// Run runs a task.
func (r *Runner) Run(id string) error {
	t := r.Task(id)

	if t == nil {
		return fmt.Errorf("task missing: %s", id)
	}

	t.ID(id)

	var e error

	select {
	case err := <-r.execAll(t):
		if err != nil {
			e = err
		}
	}

	return e
}

func (r *Runner) exec(t *task.Task) error {
	if t.UpToDate(r.ctx) {
		return errors.New("task is up to date")
	}

	backendConfig := &backendConfig.Backend{
		Log:    r.log,
		Stdin:  r.Stdin,
		Stdout: r.Stdout,
		Stderr: r.Stderr,
	}

	// Use docker if docker configuration is not nil.
	if t.Docker != nil {
		engine, err := docker.New(backendConfig)
		if err != nil {
			return err
		}
		r.engine = engine
	}

	if r.engine == nil {
		r.engine = local.New(backendConfig)
	}

	if r.verbose {
		r.log.Printf("max: using %s engine\n", r.engine.Name())
	}

	defer func() {
		r.engine.Destroy(r.ctx, t)
	}()

	if err := r.engine.Setup(r.ctx, t); err != nil {
		return err
	}

	t = r.prepareTask(t)

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

	// Prepare tasks, e.g replace arguments and environment variables.
	if err := t.Prepare(); err != nil {
		return err
	}

	// Execute task in engine.
	if err := r.engine.Exec(r.ctx, t); err != nil {
		return err
	}

	for {
		exited, err := r.engine.Wait(r.ctx, t)
		if err != nil {
			return err
		}

		if exited {
			break
		}
	}

	// Get logs from engine.
	rc, err := r.engine.Logs(r.ctx, t)
	if err != nil {
		return err
	}

	if rc != nil {
		go func() {
			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)
			log.Print(buf.String())
			rc.Close()
		}()
	}

	return nil
}

func (r *Runner) execInterval(t *task.Task) error {
	once := len(t.Interval) == 0 || r.once

	for {
		select {
		case <-r.ctx.Done():
			return errors.New("cancelled")
		default:
			break
		}

		if err := r.exec(t); err != nil {
			if once {
				err = errors.Wrap(err, "max")
				status := 1

				if strings.Contains(err.Error(), "exit status") {
					s := strings.Split(err.Error(), " ")
					if i, err := strconv.Atoi(s[len(s)-1]); err == nil {
						status = i
					}
				} else {
					log.Print(err)
				}

				os.Exit(status)
			} else {
				return err
			}
		}

		if once {
			break
		}

		// Wait until next time we should run the task.
		nextTime := cronexpr.MustParse(t.Interval).Next(time.Now())
		time.Sleep(time.Until(nextTime))
	}

	return nil
}

func (r *Runner) execAll(t *task.Task) <-chan error {
	var g errgroup.Group
	done := make(chan error)

	g.Go(func() error {
		return r.execInterval(t)
	})

	go func() {
		done <- g.Wait()
		close(done)
	}()

	return done
}

func (r *Runner) prepareTask(t *task.Task) *task.Task {
	r.parseArgs()

	if t.Variables == nil {
		t.Variables = make(map[string]string)
	}

	// Merge global variables with task variables.
	for k, v := range r.config.Variables {
		t.Variables[k] = v
	}

	t.Verbose = r.verbose

	r.args = r.config.Args

	if t.Args == nil {
		t.Args = make(map[string]interface{})
	}

	if len(r.args) > 0 {
		for k, v := range r.args {
			t.Args[k] = v
		}
	}

	t.Options(task.Log(r.log))

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

	i := 0

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
		} else if val, err := buff.ReadString(' '); err == nil || err == io.EOF {
			i++
			r.config.Variables[fmt.Sprintf("%d", i)] = strings.TrimSpace(string(rn) + val)
		}
	}
}

// Task returns a task by name if it exists.
func (r *Runner) Task(name string) *task.Task {
	if r.config == nil || r.config.Tasks[name] == nil {
		return nil
	}

	for _, n := range []string{fmt.Sprintf("%s_%s", name, runtime.GOOS), name} {
		if t := r.config.Tasks[n]; t != nil {
			return t
		}
	}

	return nil
}
