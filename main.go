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
		messageChan = make(chan string)

		// calling anonymous function that
		// closing messageChan channel and
		// set it to nil
		defer func() {
			close(messageChan)
			messageChan = nil
		}()

		// create http http.Flusher that allows
		// http handler to flush buffered data to
		// client until closed
		flusher, _ := w.(http.Flusher)
		for {
			write, err := writeData(w)
			if err != nil {
				log.Println(err)
			}
			log.Println(write)
			flusher.Flush()
		}
	}
}

// sendMessage used to write data into messageChan and flushed to client through sseStream
func sseMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Query().Get("message")
		if messageChan != nil {
			messageChan <- message
		}
	}
}

func main() {
	http.HandleFunc("/stream", sseStream())
	http.HandleFunc("/send", sseMessage())
	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}
