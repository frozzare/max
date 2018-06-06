package config

import (
	"io"
	"log"

	"github.com/frozzare/go/yaml2"
)

// Auth represents auth configuration.
type Auth struct {
	Email    string `yaml:"email"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Backend represents backend configuration.
type Backend struct {
	Log    *log.Logger
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Docker represents docker configuration.
type Docker struct {
	Auth       Auth       `yaml:"auth"`
	WorkingDir string     `yaml:"working_dir"`
	Entrypoint string     `yaml:"entrypoint"`
	Image      string     `yaml:"image"`
	Volumes    yaml2.List `yaml:"volumes"`
}
