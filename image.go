package main

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
)

func CanonicalImageName(tag string) string {
	return "docker.io/coorde/faber:" + tag
}

func PullImage(ctx context.Context, cli *client.Client, tag string) (string, error) {
	imageRef := CanonicalImageName(tag)
	reader, err := cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, reader)
	return imageRef, nil
}
