package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frozzare/max/internal/task"
	"gopkg.in/yaml.v2"
)

// Config represents a config file.
type Config struct {
	Version string
	Tasks   map[string]*task.Task
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
