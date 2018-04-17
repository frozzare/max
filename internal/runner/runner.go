package runner

import (
	"bytes"
	"io"
	"log"
	"os"
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
	args   map[string]interface{}
	All    bool
	Config *config.Config
	Once   bool
}

func (r *Runner) parseArgs() {
	input := strings.Join(os.Args[2:], " ")
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

					r.args[strings.TrimSpace(arg)] = val
				}
			}
		}
	}
}

// Run runs the tasks.
func (r *Runner) Run(id string) error {
	size := 1
	tasks := []string{id}

	// Run all tasks if defined.
	if r.All {
		if keys, err := map2.Keys(r.Config.Tasks); err == nil {
			size = len(r.Config.Tasks)
			tasks = keys.([]string)
		}
	}

	done := make(chan bool, size)

	r.args = r.Config.Args
	r.parseArgs()

	for _, k := range tasks {
		t := r.Config.Tasks[k]

		if t == nil {
			log.Fatalf("max: task missing: %s", k)
			break
		}

		once := len(t.Interval) == 0 || r.Once

		// Run task.
		go func(t *task.Task) {
			for {
				// Run deps before task.
				for _, id := range t.Deps.Values {
					dr := Runner{
						Config: r.Config,
						Once:   true,
					}

					dr.Run(id)
				}

				// Run other tasks.
				for _, id := range t.Tasks.Values {
					dr := Runner{
						Config: r.Config,
					}

					dr.Run(id)
				}

				if err := t.Run(r.args); err != nil {
					err = errors.Wrap(err, "max")

					if once {
						if !strings.Contains(err.Error(), "exit status 1") {
							log.Print(err)
						}

						os.Exit(1)
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

	return nil
}
