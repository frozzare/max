package template2

import (
	"bytes"
	"html/template"
)

var funcs = template.FuncMap{
	"isset": Isset,
}

// ExecuteString replace template tags in string with data.
func ExecuteString(c string, args interface{}) (string, error) {
	tmpl, err := template.New("main").Funcs(funcs).Parse(c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}
