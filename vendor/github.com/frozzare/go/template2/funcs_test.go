package template2

import "testing"

func TestIsset(t *testing.T) {
	var m = map[string]string{
		"foo": "",
		"bar": "foo",
	}

	if !Isset(m, "foo") {
		t.Fatal("Expected foo to be set")
	}

	if !Isset(m, "bar") {
		t.Fatal("Expected bar to be set")
	}

	if Isset(m, "foobar") {
		t.Fatal("Expected foobar not to be set")
	}

	var s = struct {
		Foo string
	}{
		Foo: "",
	}

	if !Isset(s, "Foo") {
		t.Fatal("Expected foo to be set")
	}

	if Isset(s, "foobar") {
		t.Fatal("Expected foobar not to be set")
	}

	var p = &struct {
		Foo string
	}{
		Foo: "",
	}

	if !Isset(p, "Foo") {
		t.Fatal("Expected foo to be set")
	}

	if Isset(p, "foobar") {
		t.Fatal("Expected foobar not to be set")
	}
}
