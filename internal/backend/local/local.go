package local

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/frozzare/max/internal/backend"
	"github.com/frozzare/max/internal/task"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
)

type engine struct {
}

// New creates a new local engine.
func New() backend.Engine {
	return &engine{}
}

// Setup setups local engine.
func (e *engine) Setup(task *task.Task) error {
	return nil
}

// Exec executes a task.
func (e *engine) Exec(t *task.Task) error {
	once := len(t.Interval) == 0 || false // t.Once

	for {
		if err := t.Run(); err != nil {
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

		if once {
			break
		}

		// Wait until next time we should run the task.
		nextTime := cronexpr.MustParse(t.Interval).Next(time.Now())
		time.Sleep(time.Until(nextTime))
	}

	return nil
}

// Logs returns logs from the local engine.
func (e *engine) Logs(task *task.Task) (io.ReadCloser, error) {
	return nil, nil
}

// Destroy destroys the local engine.
func (e *engine) Destroy(task *task.Task) error {
	return nil
}

// Wait check if the local engine is done or not.
func (e *engine) Wait(t *task.Task) (bool, error) {
	return true, nil
}
