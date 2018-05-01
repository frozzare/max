package backend

import (
	"io"

	"github.com/frozzare/max/internal/task"
)

// Engine defines a engine that can run tasks.
type Engine interface {
	Setup(*task.Task) error
	Exec(*task.Task) error
	Wait(*task.Task) (bool, error)
	Logs(*task.Task) (io.ReadCloser, error)
	Destroy(*task.Task) error
}
