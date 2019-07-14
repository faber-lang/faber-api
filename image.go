package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
)

func CanonicalImageName(tag string) string {
	return "docker.io/faberlang/faber:" + tag
}

func PullImage(ctx context.Context, cli *client.Client, tag string) (string, error) {
	if match, _ := regexp.MatchString(`^[\w][\w.-]{0,127}$`, tag); !match {
		return "", fmt.Errorf("invalid tag name: %s", tag)
	}

	imageRef := CanonicalImageName(tag)
	reader, err := cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, reader)
	return imageRef, nil
}
