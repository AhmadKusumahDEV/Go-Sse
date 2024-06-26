package api

import (
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

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

func Sendevent(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()
	str := "send to client"
	for client := range clients {
		client <- str
	}
}

func init() {
	app = gin.New()
	app.Use(cors.Default())

	app.GET("/addclient", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		chane := make(chan string)
		addCLient(chane)
		defer removeClient(chane)
		noOfExecution := 10
		progress := 0
		progressPercentage := float64(progress) / float64(noOfExecution) * 100

		c.SSEvent("progress", map[string]interface{}{
			"currentTask":        progress,
			"progressPercentage": progressPercentage,
			"noOftasks":          noOfExecution,
			"completed":          false,
		})
		// Flush the response to ensure the data is sent immediately
		c.Writer.Flush()

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
				c.Writer.Flush()
			case <-c.Request.Context().Done():
				return
			}
		}
	})

	app.GET("/sendevent", Sendevent)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
