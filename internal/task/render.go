package task

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/frozzare/go/env"
)

func renderEnvVariables(v string, variables map[string]string) string {
	r := regexp.MustCompile(`\$[a-zA-Z_]+[a-zA-Z0-9_]*`)
	m := r.FindAllString(v, -1)

	for _, e := range m {
		s := env.Get(e[1:])

		if len(s) == 0 {
			s = variables[e[1:]]
		}

		if len(s) == 0 {
			continue
		}

		v = strings.Replace(v, e, s, -1)
	}

	return v
}

func renderCommand(c string, args map[string]interface{}) (string, error) {
	tmpl, err := template.New("main").Parse(c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}
