package task

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/frozzare/go/env"
	"github.com/frozzare/go/structs"
	"github.com/frozzare/go/yaml2"
)

func renderEnvVariables(v string, variables map[string]string) string {
	r := regexp.MustCompile(`\$[a-zA-Z0-9_]+[a-zA-Z0-9_]*`)
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

func renderStruct(s interface{}, args map[string]interface{}, vars map[string]string) (interface{}, error) {
	fs, err := structs.Fields(s)
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		if f == nil || f.IsZero() {
			continue
		}

		skipStruct := false

		switch v := f.Value().(type) {
		case string:
			v, err = renderCommand(renderEnvVariables(v, vars), args)
			if err != nil {
				return nil, err
			}

			if err := f.Set(v); err != nil {
				return nil, err
			}
		case yaml2.List:
			skipStruct = true

			for i, k := range v.Values {
				k, err = renderCommand(renderEnvVariables(k, vars), args)
				if err != nil {
					return nil, err
				}
				v.Values[i] = k
			}

			if err := f.Set(v); err != nil {
				return nil, err
			}
		case []string:
			for i, k := range v {
				k, err = renderCommand(renderEnvVariables(k, vars), args)
				if err != nil {
					return nil, err
				}
				v[i] = k
			}

			if err := f.Set(v); err != nil {
				return nil, err
			}
		}

		if !skipStruct && (f.Kind() == reflect.Struct || f.Kind() == reflect.Ptr) {
			v, err := renderStruct(f.Value(), args, vars)
			if err != nil {
				return nil, err
			}

			if err := f.Set(v); err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}
