package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
	client, err := client.NewEnvClient()
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

// Setup setups the docker engine.
func (e *engine) Setup(ctx context.Context, t *task.Task) error {
	for _, vol := range e.volumes {
		_, err := e.client.VolumeCreate(ctx, volume.VolumesCreateBody{
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
	pullopts := types.ImagePullOptions{}

	rc, perr := e.client.ImagePull(ctx, t.Docker.Image, pullopts)
	if perr == nil {
		io.Copy(ioutil.Discard, rc)
		rc.Close()
	}

	if path, err := os.Getwd(); err == nil {
		for i, x := range t.Docker.Volumes.Values {
			t.Docker.Volumes.Values[i] = strings.Replace(x, ".:", path+":", -1)
		}
	}

	var cmds []string

	for _, c := range t.Commands.Values {
		if t.Verbose {
			cmds = append(cmds, fmt.Sprintf("echo $ %s && %s", c, c))
		} else {
			cmds = append(cmds, c)
		}
	}

	config := &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Env:          toEnv(t.Variables),
		Volumes:      toVolumes(t.Docker.Volumes.Values),
		WorkingDir:   t.Docker.WorkingDir,
		Image:        t.Docker.Image,
		Cmd:          append([]string{"sh", "-c"}, cmds...),
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

	_, err := e.client.ContainerCreate(ctx, config, hostConfig, nil, t.ID())

	if err != nil {
		return err
	}

	return e.client.ContainerStart(ctx, t.ID(), types.ContainerStartOptions{})
}

// Logs return docker logs.
func (e *engine) Logs(ctx context.Context, task *task.Task) (io.ReadCloser, error) {
	logs, err := e.client.ContainerLogs(ctx, task.ID(), types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	})

	if err != nil {
		return nil, err
	}

	rc, wc := io.Pipe()

	go func() {
		stdcopy.StdCopy(wc, wc, logs)
		logs.Close()
		wc.Close()
		rc.Close()
	}()

	return rc, nil
}

// Destroy destroys the docker container.
func (e *engine) Destroy(ctx context.Context, t *task.Task) error {
	e.client.ContainerKill(ctx, t.ID(), "9")
	e.client.ContainerRemove(ctx, t.ID(), types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})

	for _, volume := range e.volumes {
		e.client.VolumeRemove(ctx, volume.Name, true)
	}

	return nil
}

// Wait check if the conatiner is done or not.
func (e *engine) Wait(ctx context.Context, t *task.Task) (bool, error) {
	_, err := e.client.ContainerWait(ctx, t.ID())
	if err != nil {
		return false, err
	}

	info, err := e.client.ContainerInspect(ctx, t.ID())
	if err != nil {
		return false, err
	}

	if info.State.Running {
		return false, nil
	}

	return true, nil
}
