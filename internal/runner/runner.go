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
	Config  *config.Config
	Once    bool
	opts    []Option
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
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

	return r
}

// Run runs a task.
func (r *Runner) Run(id string) error {
	task := r.Task(id)

	if task == nil {
		return fmt.Errorf("task missing: %s", id)
	}

	task.SetID(id)

	var e error

	select {
	case err := <-r.execAll(task):
		if err != nil {
			e = err
		}
	}

	return e
}

func (r *Runner) exec(t *task.Task) error {
	backendConfig := &backendConfig.Backend{
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
		fmt.Println("docker")
	}

	if r.engine == nil {
		r.engine = local.New(backendConfig)
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
	once := len(t.Interval) == 0 || r.Once

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
	if t.Variables == nil {
		t.Variables = make(map[string]string)
	}

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
