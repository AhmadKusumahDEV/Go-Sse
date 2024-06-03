package main

import (
	"fmt"
	"log"
	"net/http"
)

// data variable
var messageChan chan string

func prepareHeaderForSSE(w http.ResponseWriter) {
	// prepare the header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeData(w http.ResponseWriter) (int, error) {
	// set data into response writer
	return fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
}

func sseStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// call prepareHeaderForSSE for start endpoint as SSE server
		prepareHeaderForSSE(w)
		// initialize messageChan
		fmt.Println("hei")
		messageChan = make(chan string)

		// calling anonymous function that
		// closing messageChan channel and
		// set it to nil
		defer func() {
			fmt.Println("done")
			close(messageChan)
			messageChan = nil
		}()

		// create http http.Flusher that allows
		// http handler to flush buffered data to
		// client until closed
		flusher, _ := w.(http.Flusher)
		for {
			// _, err := writeData(w)
			// if err != nil {
			// 	log.Println(err)
			// }
			<-messageChan
			// log.Println(write)
			fmt.Println("flushed")
			flusher.Flush()
		}
	}
}

// sendMessage used to write data into messageChan and flushed to client through sseStream
func sseMessage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if messageChan != nil {
			messageChan <- message
		} else {
			fmt.Println("kosong anjg")
		}
	}
}

func main() {
	http.HandleFunc("/stream", sseStream())
	http.HandleFunc("/send", sseMessage("you can put json in here as the data"))
	http.HandleFunc("/right", sseMessage("right"))
	http.HandleFunc("/left", sseMessage("left"))
	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:8080", nil))
}
