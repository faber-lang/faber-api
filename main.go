package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
	"github.com/heroku/docker-registry-client/registry"
)

type Options struct {
	Tag  string `json:"tag" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("%v", err)
	}

	url := "https://registry-1.docker.io/"
	hub, err := registry.New(url, "", "")
	if err != nil {
		log.Fatalf("%v", err)
	}

	r := gin.Default()
	r.POST("/compile", func(c *gin.Context) {
		var options Options
		if err := c.ShouldBind(&options); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := Compile(ctx, cli, options.Tag, options.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, res)
	})
	r.GET("/tags", func(c *gin.Context) {
		tags, err := hub.Tags("coorde/faber")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, tags)
	})
	r.Run()
}
