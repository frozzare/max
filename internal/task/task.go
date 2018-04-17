package task

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/frozzare/go/env"
	"github.com/frozzare/max/pkg/exec"
	"github.com/frozzare/max/pkg/log"
)

// Task represents a task.
type Task struct {
	Args     []interface{}
	Commands []string
	Deps     []string
	Interval string
	Summary  string
	Tasks    []string
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

// Run runs a task.
func (t *Task) Run(args []interface{}) error {
	if len(args) == 0 && len(t.Args) > 0 {
		args = t.Args
	}

	for _, c := range t.Commands {
		c = t.appendEnvVariables(c)

		if len(args) > 0 {
			c = fmt.Sprintf(c, args...)
		}

		res, err := exec.Cmd(c)
		if err != nil {
			log.Log(c)
			return err
		}

		if len(res) > 0 {
			log.Log(res)
		}
	}

	return nil
}
