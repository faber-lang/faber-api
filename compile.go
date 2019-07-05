package main

import (
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
)

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"code"`
}

func Compile(ctx context.Context, cli *client.Client, tag, code string) (*Result, error) {
	name, err := PullImage(ctx, cli, tag)
	if err != nil {
		return nil, err
	}

	id, err := CreateSandboxContainer(ctx, cli, name)
	if err != nil {
		return nil, err
	}

	if err := CopySourceToContainer(ctx, cli, id, code); err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	exitCode, err := cli.ContainerWait(ctx, id)
	if err != nil {
		return nil, err
	}

	stdout, stderr, err := ObtainLogs(ctx, cli, id)
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return nil, err
	}

	return &Result{Stdout: stdout, Stderr: stderr, ExitCode: int(exitCode)}, nil
}
