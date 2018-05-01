package docker

// Volume represents a Docker volume.
type Volume struct {
	Name       string
	Driver     string
	DriverOpts map[string]string
}
