package main

import (
	"log"
	"net/http"
	"io"
	"os"
	"bytes"
	"io/ioutil"
	"archive/tar"

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

func createSourceArchive(source string) (string, error) {
	file, err := ioutil.TempFile("", "faber.tar")
	if err != nil {
		return "", err
	}
	defer file.Close()

	t := tar.NewWriter(file)
	defer t.Close()

	content := []byte(source)

	header := &tar.Header{
		Name: "/faber.fab",
		Size: int64(len(content)),
	}

	if err := t.WriteHeader(header); err != nil {
		return "", err
	}

	if _, err := t.Write(content); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func copySourceToContainer(ctx context.Context, cli *client.Client, id string, code string) error {
	archive, err := createSourceArchive(code)
	if err != nil {
		return err
	}

	reader, err := os.Open(archive)
	if err != nil {
		return err
	}

	return cli.CopyToContainer(ctx, id, "/", reader, types.CopyToContainerOptions{})
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
		Cmd: []string{"cat", "/faber.fab"},
		Tty: true,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := copySourceToContainer(ctx, cli, resp.ID, code); err != nil {
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
