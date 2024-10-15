package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Message struct {
	ID      int
	Topic   string
	Content string
}

type Client struct {
	ID   int
	Conn http.ResponseWriter
	Ch   chan Message
}

var (
	topics    = make(map[string][]Client)
	messageId int
	mu        sync.Mutex
)

func main() {
	fmt.Println("Starting Infocenter Service...")
	fmt.Println("Server is running at http://localhost:8080/infocenter")

	// Registering handlers
	http.HandleFunc("/infocenter/", infocenterHandler)

	// http.HandleFunc("/infocenter/", receiveMessages) // Note: Route should differentiate based on method (GET/POST)

	http.ListenAndServe(":8080", nil)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head><title>Infocenter Service</title></head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<h1>Welcome to the Infocenter Service!</h1>")
	fmt.Fprintln(w, "<h2>Instructions:</h2>")
	fmt.Fprintln(w, "<p>To subscribe to a topic, use:</p>")
	fmt.Fprintln(w, "<pre>GET /infocenter/{topic}</pre>")
	fmt.Fprintln(w, "<p>Replace <strong>{topic}</strong> with the name of your topic.</p>")
	fmt.Fprintln(w, "<p>You can send messages to a topic using:</p>")
	fmt.Fprintln(w, "<pre>POST /infocenter/{topic}</pre>")
	fmt.Fprintln(w, "<p>with the message in the body (plain text, not JSON).</p>")
	fmt.Fprintln(w, "<h3>Example:</h3>")
	fmt.Fprintln(w, "<pre>POST /infocenter/suniukai</pre>")
	fmt.Fprintln(w, "<pre>Body: labas</pre>")
	fmt.Fprintln(w, "<p>You can test this using Postman or any HTTP client.</p>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func infocenterHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Path[len("/infocenter/"):] // Extract topic from URL

	if topic == "" {
		homepageHandler(w, r) // Call the homepage if no topic is specified
		return
	}

	if r.Method == http.MethodGet {
		// Handle subscription to messages
		// receiveMessages(w, r, topic)
	} else if r.Method == http.MethodPost {
		// Handle sending messages
		// sendMessage(w, r, topic)
	} else {
		// Handle unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
