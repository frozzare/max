package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frozzare/max/internal/task"
	"gopkg.in/yaml.v2"
)

var (
	// ErrUnmarshal is returned when config can't be unmarshaled.
	ErrUnmarshal = errors.New("max: can't unmarshal config value")
)

// Config represents a config file.
type Config struct {
	Args    map[string]interface{}
	Tasks   map[string]*task.Task
	Version string
}

type base struct {
	Args    map[string]interface{}
	Tasks   map[string]interface{}
	Version string
}

// UnmarshalYAML implements yaml packages interface to unmarshal custom values.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var b *base

	if err := unmarshal(&b); err == nil {
		c.Args = b.Args
		c.Tasks = make(map[string]*task.Task)
		c.Version = b.Version

		// Loop over tasks to include and convert existing maps to tasks.
		for k, v := range b.Tasks {
			switch r := v.(type) {
			case string:
				if content, err := ioutil.ReadFile(r); err == nil {
					var t *task.Task
					if err := yaml.Unmarshal([]byte(content), &t); err == nil {
						c.Tasks[k] = t
					} else {
						return ErrUnmarshal
					}
				}
			case map[interface{}]interface{}:
				var t *task.Task

				if buf, err := yaml.Marshal(r); err == nil {
					if err := yaml.Unmarshal(buf, &t); err == nil {
						c.Tasks[k] = t
					} else {
						return ErrUnmarshal
					}
				} else {
					return ErrUnmarshal
				}
			}
		}

		return nil
	}

	return ErrUnmarshal
}

func convert(a interface{}) interface{} {
	switch v := a.(type) {
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range v {
			res[k.(string)] = v
		}
		return res
	default:
		return v
	}
}

// ReadContent creates a new config struct from a string.
func ReadContent(content string) (*Config, error) {
	var config *Config

	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return &Config{}, err
	}

	return config, nil
}

// ReadFile creates a new config struct from a yaml file.
func ReadFile(args ...string) (*Config, error) {
	var file string
	var path string
	var err error

	if len(args) > 0 && args[0] != "" {
		if _, err := os.Stat(args[0]); err == nil {
			file = args[0]
			path = filepath.Dir(file)
		} else {
			path = args[0]
		}
	}

	if !strings.HasPrefix("/", path) {
		path, err = os.Getwd()
	}

	files := []string{fmt.Sprintf("max_%s.yml", runtime.GOOS), "max.yml"}
	if len(file) > 0 {
		files = append([]string{file}, files...)
	}

	var dat []byte
	for _, name := range files {
		if len(dat) > 0 {
			break
		}

		file := filepath.Join(path, name)
		dat, err = ioutil.ReadFile(file)
	}

	if err != nil {
		return nil, err
	}

	var config *Config

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return &Config{}, err
	}

	return config, nil
}
