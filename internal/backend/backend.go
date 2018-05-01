package backend

import (
	"context"
	"io"

	"github.com/frozzare/max/internal/task"
)

// Engine defines a engine that can run tasks.
type Engine interface {
	Destroy(context.Context, *task.Task) error
	Exec(context.Context, *task.Task) error
	Logs(context.Context, *task.Task) (io.ReadCloser, error)
	Name() string
	Setup(context.Context, *task.Task) error
	Wait(context.Context, *task.Task) (bool, error)
}
