package main

import (
	"embed"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var f embed.FS

var clients map[chan string]struct{}
var lock sync.RWMutex

func addCLient(clien chan string) {
	lock.Lock()
	defer lock.Unlock()
	clients[clien] = struct{}{}
}

func removeClient(clien chan string) {
	lock.Lock()
	defer lock.Unlock()
	delete(clients, clien)
	close(clien)
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// SSE endpoint
	router.GET("/", func(ctx *gin.Context) {
		num = 1
	})

	router.GET("/sendevent", Sendevent)
	router.GET("/progress", progressor)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

var num int = 0

func Sendevent(c *gin.Context) {
	str := "data: " + string(num)
	chane <- str
}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
// 		c.Writer.Header().Set("Content-Type", "text/event-stream")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")
// 		c.Writer.Header().Set("Connection", "keep-alive")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

func progressor(c *gin.Context) {
	chane := make(chan string)

	noOfExecution := 10
	progress := 0
	progressPercentage := float64(progress) / float64(noOfExecution) * 100

	c.SSEvent("progress", map[string]interface{}{
		"currentTask":        progress,
		"progressPercentage": progressPercentage,
		"noOftasks":          noOfExecution,
		"completed":          false,
	})

	defer func() {
		fmt.Println("done")
	}()
	// Flush the response to ensure the data is sent immediately
	c.Writer.Flush()

	for {
		<-chane
		c.SSEvent("progress", map[string]interface{}{
			"currentTask":        progress,
			"progressPercentage": progressPercentage,
			"noOftasks":          noOfExecution,
			"completed":          false,
		})
		// Flush the response to ensure the data is sent immediately
		c.Writer.Flush()
	}
}
