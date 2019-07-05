package main

import (
	"bytes"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	units "github.com/docker/go-units"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
)

func CopySourceToContainer(ctx context.Context, cli *client.Client, id string, code string) error {
	archive, err := CreateSourceArchive(code)
	if err != nil {
		return err
	}

	reader, err := os.Open(archive)
	if err != nil {
		return err
	}

	return cli.CopyToContainer(ctx, id, "/", reader, types.CopyToContainerOptions{})
}

func ObtainLogs(ctx context.Context, cli *client.Client, id string) (string, string, error) {
	stdout, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", "", err
	}
	stdoutbuf := new(bytes.Buffer)
	stdoutbuf.ReadFrom(stdout)

	stderr, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		return "", "", err
	}
	stderrbuf := new(bytes.Buffer)
	stderrbuf.ReadFrom(stderr)

	return stdoutbuf.String(), stderrbuf.String(), nil
}

func CreateSandboxContainer(ctx context.Context, cli *client.Client, name string) (string, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: name,
		Cmd:   []string{"fabrun", "/faber.fab"},
		Tty:   true,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:    50000000,
			CPUPeriod: 1000000,
			CPUQuota:  200000,
			Ulimits: []*units.Ulimit{&units.Ulimit{
				Name: "cpu",
				Soft: 1,
				Hard: 1,
			}},
		},
		Privileged:  false,
		NetworkMode: "none",
	}, nil, "")
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
