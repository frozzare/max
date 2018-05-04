package config

import (
	"io"
	"log"
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
	Context    string
	Entrypoint string
	Image      string
	Volumes    []string
}
