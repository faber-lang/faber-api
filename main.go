package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
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
	r.Run()
}
