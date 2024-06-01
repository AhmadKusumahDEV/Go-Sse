package main

import (
	"embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


var f embed.FS

func main() {
	router := gin.Default()

	// SSE endpoint
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tpl", nil)
	})
	router.GET("/progress", progressor)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func progressor(c *gin.Context) {
	noOfExecution := 10
	progress := 0
	for progress <= noOfExecution {
		progressPercentage := float64(progress) / float64(noOfExecution) * 100

		c.SSEvent("progress", map[string]interface{}{
			"currentTask":        progress,
			"progressPercentage": progressPercentage,
			"noOftasks":          noOfExecution,
			"completed":          false,
		})
		// Flush the response to ensure the data is sent immediately
		c.Writer.Flush()

		progress += 1
		time.Sleep(2 * time.Second)
	}

	c.SSEvent("progress", map[string]interface{}{
		"completed":          true,
		"progressPercentage": 100,
	})

	// Flush the response to ensure the data is sent immediately
	c.Writer.Flush()

}