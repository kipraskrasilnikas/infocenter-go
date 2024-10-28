package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	ID      int
	Topic   string
	Content string
}

var (
	mu        sync.Mutex
	topics    = make(map[string]chan Message)
	messageId int
)

func main() {
	fmt.Println("Starting Infocenter Service...")
	fmt.Println("Server is running at http://localhost:8080/infocenter")

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Registering handlers
	http.HandleFunc("/infocenter/", infocenterHandler)
	http.ListenAndServe(":8080", nil)
}

func infocenterHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Path[len("/infocenter/"):] // Extract topic from URL

	if topic == "" {
		homepageHandler(w, r) // Call the homepage if no topic is specified
		return
	}

	if r.Method == http.MethodGet {
		receiveMessages(w, topic)
	} else if r.Method == http.MethodPost {
		sendMessage(w, r, topic)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// TO DO: Remake this to handle multiple clients
func receiveMessages(w http.ResponseWriter, topic string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	mu.Lock()
	if _, exists := topics[topic]; !exists {
		topics[topic] = make(chan Message)
	}
	mu.Unlock()

	fmt.Printf("Listening to topic: %s\n", topic)
	ch := topics[topic]

	timeoutTime := 30
	timeout := time.After(time.Duration(timeoutTime) * time.Second)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return // Channel is closed, exit
			}

			fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Content)
			w.(http.Flusher).Flush()
		case <-timeout:
			// Send a timeout event before disconnecting
			fmt.Fprintf(w, "id: %d\nevent: timeout\ndata: %ds\n\n", messageId, timeoutTime)
			w.(http.Flusher).Flush() // Flush the response writer

			cleanupTopic(topic)
			return
		}
	}
}

// Potential improvement: save messages when clients are not listening and let clients catch up after they start
func sendMessage(w http.ResponseWriter, r *http.Request, topic string) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()

	if ch, exists := topics[topic]; exists {
		messageId++ // Increment the message ID
		message := Message{ID: messageId, Topic: topic, Content: string(body)}

		// Send the message to the channel for this topic
		select {
		case ch <- message:
			// Message sent successfully
		default:
			// If the channel is full, skip send to avoid blocking
			fmt.Printf("Message dropped for topic: %s\n", topic)
		}
	}

	mu.Unlock()
	w.WriteHeader(http.StatusNoContent) // Send HTTP 204 No Content response
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html") // Serve the static HTML file
}

func cleanupTopic(topic string) {
	mu.Lock()
	defer mu.Unlock()
	if ch, exists := topics[topic]; exists {
		close(ch)
		delete(topics, topic)
	}
}
