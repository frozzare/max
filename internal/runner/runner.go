package runner

import (
	"time"

	"github.com/frozzare/go/map2"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
	"github.com/frozzare/max/pkg/log"
	golog "github.com/go-log/log"
	"github.com/gorhill/cronexpr"
)

// Runner represents a the runner.
type Runner struct {
	All    bool
	Config *config.Config
	Once   bool
	Quiet  bool
}

func toInterface(list []string) []interface{} {
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}
	return vals
}

// Run runs the tasks.
func (r *Runner) Run(id string, args ...string) error {
	// Output help usage if requested.
	if id == "help" && len(args) == 1 {
		id = args[0]
		t := r.Config.Tasks[id]

		if t == nil {
			log.Fatalf("Task missing: %s", id)
		}

		if len(t.Usage) == 0 {
			return nil
		}

		log.Logf("Usage:\n  max %s %s", id, t.Usage)
		log.Logf("\nSummary:\n  %s", t.Summary)

		return nil
	}

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
	errs := make(chan error)

	// Set output log to quiet.
	if r.Quiet {
		log.SetLogger(golog.DefaultLogger)
	}

	for _, k := range tasks {
		t := r.Config.Tasks[k]

		if t == nil {
			log.Fatalf("Task missing: %s", k)
			break
		}

		// Run task.
		go func(t *task.Task) {
			for {
				// Run deps before task.
				for _, id := range t.Deps {
					dr := Runner{
						Config: r.Config,
						Once:   true, // deps can only run once since the are runned in another tasks interval
						Quiet:  r.Quiet,
					}

					dr.Run(id, args...)
				}

				// Run other tasks.
				for _, id := range t.Tasks {
					dr := Runner{
						Config: r.Config,
						Quiet:  r.Quiet,
					}

					dr.Run(id, args...)
				}

				if err := t.Run(toInterface(args)); err != nil {
					log.Log(err)
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
			close(errs)
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
