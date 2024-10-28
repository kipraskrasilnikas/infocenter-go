# InfoCenter

A project that allows clients to achieve almost real-time communication between each other by sending messages, using concurrency to send or receive messages at the same time.

## How to run the application

- Clone the Git repository to your local machine using the command `git clone https://github.com/kipraskrasilnikas/infocenter-go/`
- Open a terminal or command prompt and navigate to the directory where you cloned the repository.
- Run the command `go build -o infocenter.exe` to build the executable file.
- Run the command `./infocenter.exe` to launch the project. This will start the server on port 8080.
- Open a web browser and navigate to `http://localhost:8080/infocenter/`. This will show the home page of the project.

## How the project works

- The `main.go` file contains the entry point for the program. It starts an HTTP server, registers an HTTP request handler for the `/infocenter/` endpoint.
- The homepage is shown to the user with the topic and message input fields, subscribe and send message buttons and a box for incoming messages.
- When a client inputs a topic and clicks the "Subscribe to topic" button, a `GET` request is made to the `/infocenter/<topic>` endpoint. This creates a channel for the topic in the topics map and opens a persistent connection, which listens for events in text/event-stream format.
- When a client adds a message, topic and clicks "Send message", a `POST` request is made to the `/infocenter/<topic>` endpoint. This creates an incremented message object, sends it to the corresponding channel (topic) and returns HTTP status code 204.
- If the topic is subscribed to, the message is sent as an SSE (Server-sent event) to the `GET` request response. The message is shown to the subscribed user.
- The subscription times out after 30 seconds. The topic is deleted from the topics map and the channel is closed.
- All `Message` object and `topics` map operations are protected by a `sync.Mutex` to allow concurrent read access and exclusive write access.
