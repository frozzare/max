package task

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/imdario/mergo"

	"github.com/frozzare/go/env"
	"github.com/frozzare/max/pkg/exec"
	"github.com/frozzare/max/pkg/yamllist"
)

// Task represents a task.
type Task struct {
	Args     map[string]interface{}
	Commands yamllist.List
	Deps     []string
	Dir      string
	Interval string
	Summary  string
	Tasks    yamllist.List
	Usage    string
	Verbose  bool
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

func (t *Task) appendEnvVariables(v string) string {
	r := regexp.MustCompile(`\$[a-zA-Z_]+[a-zA-Z0-9_]*`)
	m := r.FindAllString(v, -1)

	for _, e := range m {
		v = strings.Replace(v, e, env.Get(e[1:]), -1)
	}

	return v
}

func (t *Task) appendArguments(c string) (string, error) {
	tmpl, err := template.New("main").Parse(c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t.Args); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *Task) prepareString(c string) (string, error) {
	return t.appendArguments(t.appendEnvVariables(c))
}

// Run runs task commands.
func (t *Task) Run(args map[string]interface{}) error {
	if len(args) > 0 {
		if err := mergo.Merge(&args, t.Args); err != nil {
			return err
		}

		t.Args = args
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

		// Execute command.
		if err := exec.Exec(c, t.Dir); err != nil {
			log.Print(c)
			return err
		}
	}

	return nil
}
