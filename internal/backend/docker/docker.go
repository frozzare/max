package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/frozzare/max/internal/backend"
	"github.com/frozzare/max/internal/backend/config"
	"github.com/frozzare/max/internal/task"
)

type engine struct {
	config  *config.Backend
	client  client.APIClient
	volumes []Volume
}

// New returns a new Docker Engine using the given client.
func New(config *config.Backend) (backend.Engine, error) {
	client, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &engine{
		config: config,
		client: client,
		volumes: []Volume{
			{
				Name:   "max_default",
				Driver: "local",
			},
		},
	}, nil
}

// Name returns engine name.
func (e *engine) Name() string {
	return "docker"
}

func (e *engine) containerName(t *task.Task) string {
	return fmt.Sprintf("max_%s", t.ID())
}

// Setup setups the docker engine.
func (e *engine) Setup(ctx context.Context, t *task.Task) error {
	for _, vol := range e.volumes {
		_, err := e.client.VolumeCreate(ctx, volume.CreateOptions{
			Name:       vol.Name,
			Driver:     vol.Driver,
			DriverOpts: vol.DriverOpts,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Exec execute a task in a docker container.
func (e *engine) Exec(ctx context.Context, t *task.Task) error {
	pullopts := image.PullOptions{}

	// Add authentication credentials if any.
	if t.Docker.Auth != nil && len(t.Docker.Auth.Username) > 0 && len(t.Docker.Auth.Password) > 0 {
		b, err := json.Marshal(t.Docker.Auth)
		if err != nil {
			return err
		}
		pullopts.RegistryAuth = base64.URLEncoding.EncodeToString(b)
	}

	rc, err := e.client.ImagePull(ctx, t.Docker.Image, pullopts)
	if err == nil {
		io.Copy(io.Discard, rc)
		rc.Close()
	}

	if err != nil {
		return err
	}

	if path, err := os.Getwd(); err == nil {
		for i, x := range t.Docker.Volumes.Values {
			t.Docker.Volumes.Values[i] = strings.Replace(x, ".:", path+":", -1)
		}
	}

	var cmds []string

	for _, c := range t.Commands.Values {
		if e.config.Verbose {
			cmds = append(cmds, fmt.Sprintf("echo $ %s && %s", c, c))
		} else {
			cmds = append(cmds, c)
		}
	}

	cmd := generateCommand(generateScript(cmds))

	config := &container.Config{
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env:          toEnv(t.Variables),
		Volumes:      toVolumes(t.Docker.Volumes.Values),
		WorkingDir:   t.Docker.WorkingDir,
		Image:        t.Docker.Image,
		Cmd:          append([]string{"sh", "-c"}, cmd),
		Entrypoint:   strings.Split(t.Docker.Entrypoint, " "),
	}

	if len(config.Entrypoint) == 0 {
		config.Entrypoint = []string{"/bin/sh", "-c"}
	}

	if len(config.WorkingDir) == 0 {
		config.WorkingDir = t.Dir
	}

	hostConfig := &container.HostConfig{
		Binds: t.Docker.Volumes.Values,
	}

	networkConfig := &network.NetworkingConfig{}

	_, err = e.client.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, e.containerName(t),)

	if err != nil {
		return err
	}

	return e.client.ContainerStart(ctx, e.containerName(t), container.StartOptions{})
}

// Logs return docker logs.
func (e *engine) Logs(ctx context.Context, t *task.Task) (io.ReadCloser, error) {
	return e.client.ContainerLogs(ctx, e.containerName(t), container.LogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	})
}

// Destroy destroys the docker container.
func (e *engine) Destroy(ctx context.Context, t *task.Task) error {
	e.client.ContainerKill(ctx, e.containerName(t), "9")
	e.client.ContainerRemove(ctx, e.containerName(t), container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})

	for _, volume := range e.volumes {
		e.client.VolumeRemove(ctx, volume.Name, true)
	}

	return nil
}

// Wait check if the container is done or not.
func (e *engine) Wait(ctx context.Context, t *task.Task) (bool, error) {
	_, errC := e.client.ContainerWait(ctx, e.containerName(t), container.WaitConditionNextExit)
	if err := <-errC; err != nil {
		return false, err
	}

	info, err := e.client.ContainerInspect(ctx, e.containerName(t))
	if err != nil {
		return false, err
	}

	return !info.State.Running, nil
}
