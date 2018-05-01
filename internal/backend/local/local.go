package local

import (
	"context"
	"io"

	"github.com/frozzare/max/internal/backend"
	"github.com/frozzare/max/internal/task"
)

type engine struct {
}

// New creates a new local engine.
func New() backend.Engine {
	return &engine{}
}

// Setup setups local engine.
func (e *engine) Setup(ctx context.Context, t *task.Task) error {
	return nil
}

// Exec executes a task.
func (e *engine) Exec(ctx context.Context, t *task.Task) error {
	return t.Run()
}

// Logs returns logs from the local engine.
func (e *engine) Logs(ctx context.Context, t *task.Task) (io.ReadCloser, error) {
	return nil, nil
}

// Destroy destroys the local engine.
func (e *engine) Destroy(ctx context.Context, t *task.Task) error {
	return nil
}

// Wait check if the local engine is done or not.
func (e *engine) Wait(ctx context.Context, t *task.Task) (bool, error) {
	return true, nil
}
