package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/moby/moby/client"
	"golang.org/x/net/context"
)

type Options struct {
	Tag  string `json:"tag" binding:"required"`
	Code string `json:"code" binding:"required"`
	Save bool   `json:"save"`
}

type Return struct {
	Result
	ID string `json:"id"`
}

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	url := "https://registry-1.docker.io/"
	hub, err := registry.New(url, "", "")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	db, err := InitDB()
	if err != nil {
		log.Fatalf("%v", err)
		return
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
		ret := Return{Result: *res}

		if options.Save {
			id, err := Save(db, options, *res)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ret.ID = id

		}
		c.JSON(200, ret)
	})

	r.GET("/restore/:id", func(c *gin.Context) {
		id := c.Param("id")
		opts, res, err := Restore(db, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"options": opts,
			"result":  res,
		})
	})

	r.GET("/tags", func(c *gin.Context) {
		tags, err := hub.Tags("coorde/faber")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, tags)
	})

	if domain := os.Getenv("FABER_API_AUTOTLS_DOMAIN"); domain != "" {
		log.Fatal(autotls.Run(r, domain))
	} else {
		r.Run()
	}
}
