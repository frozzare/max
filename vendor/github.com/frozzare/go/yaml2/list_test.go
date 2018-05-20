package yaml2

import (
	"io/ioutil"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type task struct {
	Commands List
}

type test struct {
	Tasks map[string]*task
}

func TestList(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/list.yml")

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	var s *test

	if err := yaml.Unmarshal([]byte(content), &s); err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	if !reflect.DeepEqual(s.Tasks["list"].Commands.Values, []string{"List item"}) {
		t.Errorf("Expected: [List item], got: %v", s.Tasks["list"].Commands.Values)
	}

	if !reflect.DeepEqual(s.Tasks["list2"].Commands.Values, []string{"List item"}) {
		t.Errorf("Expected: [List item], got: %v", s.Tasks["list2"].Commands.Values)
	}
}
