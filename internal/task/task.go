package task

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"text/template"

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

// Run runs task commands.
func (t *Task) Run(args map[string]interface{}) error {
	if len(args) > 0 {
		t.Args = args
	}

	for _, c := range t.Commands.Values {
		c = t.appendEnvVariables(c)

		c, err := t.appendArguments(c)
		if err != nil {
			return err
		}

		if err := exec.Exec(c, t.Dir); err != nil {
			log.Print(c)
			return err
		}
	}

	return nil
}
