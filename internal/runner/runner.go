package runner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/frozzare/go/map2"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
)

// Runner represents a the runner.
type Runner struct {
	args    map[string]interface{}
	All     bool
	Config  *config.Config
	Once    bool
	Verbose bool
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

// Run runs the tasks.
func (r *Runner) Run(id string) {
	size := 1
	tasks := []string{id}

	// Run all tasks if defined.
	if r.All {
		keys, err := map2.Keys(r.Config.Tasks)
		if err != nil {
			log.Fatalf("max: %s", err.Error())
			return
		}

		size = len(r.Config.Tasks)
		tasks = keys.([]string)
	}

	done := make(chan bool, size)

	r.args = r.Config.Args
	r.parseArgs()

	for _, k := range tasks {
		t := r.Task(k)

		if t == nil {
			log.Fatalf("max: task missing: %s", k)
			break
		}

		once := len(t.Interval) == 0 || r.Once

		// Run task.
		go func(t *task.Task) {
			for {
				// Run deps before task.
				for _, id := range t.Deps {
					dr := Runner{
						Config:  r.Config,
						Once:    true,
						Verbose: r.Verbose,
					}

					dr.Run(id)
				}

				// Run other tasks.
				for _, id := range t.Tasks.Values {
					dr := Runner{
						Config:  r.Config,
						Verbose: r.Verbose,
					}

					dr.Run(id)
				}

				// Merge global variables with task variables.
				for k, v := range r.Config.Variables {
					t.Variables[k] = v
				}

				t.Verbose = r.Verbose

				if err := t.Run(r.args); err != nil {
					err = errors.Wrap(err, "max")

					if once {
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
						log.Print(err)
					}
				}

				// If no internal or only once flag is used we should break it.
				if once {
					break
				}

				// Wait until next time we should run the task.
				nextTime := cronexpr.MustParse(t.Interval).Next(time.Now())
				time.Sleep(time.Until(nextTime))
			}

			done <- true
		}(t)
	}

	// Wait for all tasks to be done.
	for {
		if len(done) == size {
			close(done)
			break
		}

		time.Sleep(1 * time.Second)
	}
}
