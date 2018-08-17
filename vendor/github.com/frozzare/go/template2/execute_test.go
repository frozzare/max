package template2

import "testing"

func TestExecuteString(t *testing.T) {
	s, err := ExecuteString("<p>{{ .name }}</p>", map[string]string{"name": "Fredrik"})
	if err != nil {
		t.Fatal(err)
	}

	if s != "<p>Fredrik</p>" {
		t.Fatal("Expected string to be equal")
	}
}
