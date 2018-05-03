package task

import (
	"log"
	"strings"

	"github.com/frozzare/max/internal/backend/config"
	"github.com/frozzare/max/pkg/yamllist"
)

// Task represents a task.
type Task struct {
	Args      map[string]interface{}
	Commands  yamllist.List
	Deps      []string
	Dir       string
	Docker    *config.Docker
	id        string
	Interval  string
	Summary   string
	Tasks     yamllist.List
	Usage     string
	Variables map[string]string
	Verbose   bool
}

// ID returns the task id.
func (t *Task) ID() string {
	return t.id
}

// SetID sets the task id.
func (t *Task) SetID(id string) {
	t.id = id
}

// PrintUsage print usage of task.
func (t *Task) PrintUsage(id string) {
	if len(t.Usage) != 0 {
		log.Printf("Usage:\n  max %s %s\n", id, t.Usage)
	}

	if len(t.Summary) != 0 {
		log.Printf("Summary:\n  %s", t.Summary)
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

	for i, c := range t.Commands.Values {
		c, err := t.prepareString(c)
		if err != nil {
			return err
		}

		t.Commands.Values[i] = c
	}

	return nil
}

func (t *Task) prepareString(c string) (string, error) {
	return renderCommand(renderEnvVariables(c, t.Variables), t.Args)
}
