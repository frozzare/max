package task

import (
	"log"
	"strings"

	"github.com/frozzare/max/pkg/exec"
	"github.com/frozzare/max/pkg/yamllist"
)

// Task represents a task.
type Task struct {
	Args      map[string]interface{}
	Commands  yamllist.List
	Deps      []string
	Dir       string
	Interval  string
	Summary   string
	Tasks     yamllist.List
	Usage     string
	Variables map[string]string
	Verbose   bool
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

func (t *Task) prepareString(c string) (string, error) {
	return renderCommand(renderEnvVariables(c, t.Variables), t.Args)
}

// Run runs task commands.
func (t *Task) Run(args map[string]interface{}) error {
	if t.Args == nil {
		t.Args = make(map[string]interface{})
	}

	if len(args) > 0 {
		for k, v := range args {
			t.Args[k] = v
		}
	}

	// Support usage of environment variabels and arguments in directory field.
	d, err := t.prepareString(t.Dir)
	if err != nil {
		return err
	}

	// Trim spaces if any exists.
	t.Dir = strings.TrimSpace(d)

	for _, c := range t.Commands.Values {
		// Prepare string with environment variables and arguments.
		c, err := t.prepareString(c)
		if err != nil {
			return err
		}

		if t.Verbose {
			log.Print(c)
		}

		opts := &exec.Options{
			Dir:     t.Dir,
			Env:     toEnv(t.Variables),
			Command: c,
		}

		// Execute command.
		if err := exec.Exec(opts); err != nil {
			log.Print(c)
			return err
		}
	}

	return nil
}
