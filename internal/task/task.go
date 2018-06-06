package task

import (
	"context"
	"log"
	"os"

	"github.com/frozzare/go/yaml2"
	"github.com/frozzare/max/internal/backend/config"
	"github.com/frozzare/max/internal/exec"
)

// Task represents a task.
type Task struct {
	Args      map[string]interface{}
	Commands  yaml2.List
	Deps      []string
	Dir       string
	Docker    *config.Docker
	Interval  string
	Summary   string
	Status    yaml2.List
	Tasks     yaml2.List
	Usage     string
	Variables map[string]string

	id  string      `structs:"-"`
	log *log.Logger `structs:"-"`
}

// ID returns the task id.
func (t *Task) ID(id ...string) string {
	if len(id) > 0 {
		t.id = id[0]
	}

	return t.id
}

// Options sets task options.
func (t *Task) Options(opts ...Option) {
	for _, opts := range opts {
		opts(t)
	}
}

// PrintUsage print usage of task.
func (t *Task) PrintUsage(id string) {
	if t.log == nil {
		t.log = log.New(os.Stderr, "", 0)
	}

	t.log.Println()

	if len(t.Usage) != 0 {
		t.log.Printf("Usage:\n\n  max %s %s\n\n", id, t.Usage)
	}

	if len(t.Summary) != 0 {
		t.log.Printf("Summary:\n\n  %s", t.Summary)
	}
}

// Prepare prepares the command and directory.
func (t *Task) Prepare() error {
	v, err := renderStruct(t, t.Args, t.Variables)
	if err != nil {
		return err
	}

	t = v.(*Task)

	return nil
}

// UpToDate determine if a task is up to date using status commands.
func (t *Task) UpToDate(ctx context.Context) bool {
	if len(t.Status.Values) == 0 {
		return false
	}

	for _, c := range t.Status.Values {
		opts := &exec.Options{
			Context: ctx,
			Dir:     t.Dir,
			Env:     toEnv(t.Variables),
			Command: c,
		}

		if err := exec.Exec(opts); err != nil {
			return false
		}
	}

	return true
}
