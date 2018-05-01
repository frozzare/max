package config

// Docker represents docker configuration.
type Docker struct {
	Context    string
	Entrypoint string
	Image      string
	Volumes    []string
}
