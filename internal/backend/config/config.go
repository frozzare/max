package config

import "io"

// Backend represents backend configuration.
type Backend struct {
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
