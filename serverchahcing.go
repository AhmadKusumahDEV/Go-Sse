package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseCahcing struct {
	Status int `json:"code"`
	Dataaa any `json:"data"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, ETag, X-Custom-Header-1, X-Custom-Header-2, X-Custom-Header-3")
		// c.Writer.Header().Set("Access-Control-Expose-Headers", "ETag")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/caching", caching)
	router.GET("/resource", func(c *gin.Context) {
		c.Header("X-Custom-Header-1", "Value1")
		c.Header("X-Custom-Header-2", "Value2")
		c.Header("X-Custom-Header-3", "Value3")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	router.Run(":8080")
}

var etag string = "some-unique-string"

func caching(ctx *gin.Context) {
	res := responseCahcing{
		Status: 200,
		Dataaa: "localstorage",
	}
	ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, ETag, X-Custom-Header-1, X-Custom-Header-2, X-Custom-Header-3")
	ctx.Header("Access-Control-Expose-Headers", "ETag")
	ctx.Header("Cache-Control", "public, max-age=3600")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("ETag", etag)
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
	ctx.Header("Last-Modified", "Mon, 02 Jan 2006 15:04:05 MST")
	ctx.Header("age", "3600")
	ctx.JSON(http.StatusAccepted, res)
}
