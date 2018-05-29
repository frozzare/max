package config

import (
	"io"
	"log"

	"github.com/frozzare/max/pkg/yamllist"
)

// Backend represents backend configuration.
type Backend struct {
	Log    *log.Logger
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Docker represents docker configuration.
type Docker struct {
	WorkingDir string        `yaml:"working_dir"`
	Entrypoint string        `yaml:"entrypoint"`
	Image      string        `yaml:"image"`
	Volumes    yamllist.List `yaml:"volumes"`
}
