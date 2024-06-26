package main

import (
	"embed"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var f embed.FS

var clients map[chan string]struct{} = make(map[chan string]struct{})
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
	lock.RLock()
	defer lock.RUnlock()
	str := "data: " + string(num)
	for client := range clients {
		client <- str
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func progressor(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	chane := make(chan string)
	fmt.Println(len(clients))
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported!"})
	}
	addCLient(chane)
	defer removeClient(chane)
	fmt.Println(clients)
	noOfExecution := 10
	progress := 0
	progressPercentage := float64(progress) / float64(noOfExecution) * 100

	c.SSEvent("progress", map[string]interface{}{
		"currentTask":        progress,
		"progressPercentage": progressPercentage,
		"noOftasks":          noOfExecution,
		"completed":          false,
	})
	flusher.Flush()
	// Flush the response to ensure the data is sent immediately
	for {
		select {
		case message := <-chane:
			c.SSEvent("progress", map[string]interface{}{
				"currentTask":        progress,
				"progressPercentage": progressPercentage,
				"noOftasks":          noOfExecution,
				"completed":          false,
				"message":            message,
			})
			flusher.Flush()
		case <-c.Request.Context().Done():
			fmt.Println("done")
			return
		}
	}
}
