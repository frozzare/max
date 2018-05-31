package task

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/frozzare/max/internal/backend/config"
	"github.com/frozzare/max/pkg/exec"
	"github.com/frozzare/max/pkg/yamllist"
)

// Task represents a task.
type Task struct {
	Args      map[string]interface{}
	Commands  yamllist.List
	Deps      []string
	Dir       string
	Docker    *config.Docker
	Interval  string
	Summary   string
	Status    yamllist.List
	Tasks     yamllist.List
	Usage     string
	Variables map[string]string

	id      string
	log     *log.Logger
	Verbose bool
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
	// Support usage of environment variabels and arguments in directory field.
	d, err := t.prepareString(t.Dir)
	if err != nil {
		return err
	}

	// Trim spaces if any exists.
	t.Dir = strings.TrimSpace(d)

	// Prepare status commands.
	cmds, err := t.prepareSlice(t.Status.Values)
	if err != nil {
		return err
	}
	t.Status.Values = cmds

	// Prepare command values.
	cmds, err = t.prepareSlice(t.Commands.Values)
	if err != nil {
		return err
	}
	t.Commands.Values = cmds

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

func (t *Task) prepareSlice(s []string) ([]string, error) {
	for i, c := range s {
		c, err := t.prepareString(c)
		if err != nil {
			return []string{}, err
		}

		s[i] = c
	}

	return s, nil
}

func (t *Task) prepareString(c string) (string, error) {
	return renderCommand(renderEnvVariables(c, t.Variables), t.Args)
}
