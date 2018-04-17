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
	All    bool
	Config *config.Config
	Once   bool
}

func (r *Runner) parseArgs() map[string]interface{} {
	args := map[string]interface{}{}
	input := strings.Join(os.Args[2:], " ")
	buff := bytes.NewBufferString(input)

	for {
		r, _, err := buff.ReadRune()

		if err != nil {
			break
		}

		if buff.Len() == 0 {
			break
		}

		if r == '-' {
			if arg, err := buff.ReadString(' '); err == nil {
				if val, err := buff.ReadString(' '); err == nil || err == io.EOF {
					if arg[0] == '-' {
						arg = arg[1:]
					}

					args[strings.TrimSpace(arg)] = val
				}
			}
		}
	}

	return args
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
	args := r.parseArgs()

	for _, k := range tasks {
		t := r.Config.Tasks[k]

		if t == nil {
			log.Fatalf("max: task missing: %s", k)
			break
		}

		// Run task.
		go func(t *task.Task) {
			for {
				// Run deps before task.
				for _, id := range t.Deps {
					dr := Runner{
						Config: r.Config,
						Once:   true,
					}

					dr.Run(id)
				}

				// Run other tasks.
				for _, id := range t.Tasks {
					dr := Runner{
						Config: r.Config,
					}

					dr.Run(id)
				}

				if err := t.Run(args); err != nil {
					log.Print(errors.Wrap(err, "max"))
				}

				// If no internal or only once flag is used we should break it.
				if len(t.Interval) == 0 || r.Once {
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
