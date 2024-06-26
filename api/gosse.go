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

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"enkripsi: get/post": "https://stockis.vercel.app/api/en",
			"dekripsi: get/post": "https://stockis.vercel.app/api/des",
		})
		return
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
