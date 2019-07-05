package main

import (
	"log"
	"net/http"
	"io"
	"os"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/moby/moby/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

type Options struct {
	Tag  string `json:"tag" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"code"`
}

func compile(ctx context.Context, cli *client.Client, tag, code string) (*Result, error) {
	imageRef := "docker.io/coorde/faber:" + tag
	reader, err := cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		Cmd: []string{"echo", "hello world"},
		Tty: true,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	exitCode, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return nil, err
	}

	stdout, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}
	stdoutbuf := new(bytes.Buffer)
	stdoutbuf.ReadFrom(stdout)

	stderr, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		return nil, err
	}
	stderrbuf := new(bytes.Buffer)
	stderrbuf.ReadFrom(stderr)

	return &Result{Stdout: stdoutbuf.String(), Stderr: stderrbuf.String(), ExitCode: int(exitCode)}, nil
}

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("%v", err)
	}

	r := gin.Default()
	r.POST("/compile", func(c *gin.Context) {
		var options Options
		if err := c.ShouldBind(&options); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		res, err := compile(ctx, cli, options.Tag, options.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, res)
	})
	r.Run()
}
