package main

import (
	"fmt"
	"net/http"
	"mygoproject/pkg/websocket"
)

// note: a ResponseWriter is an http response. When we write to it, we are sending response to client.
// note: a Request is an http request. It is the request coming from the client.
// note: our custom pool struct can register/unregister a client, contains the map of clients, and broadcast method to send a message to all clients
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	// upgrade a http protocol to a websocket protocol
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	// create instance of custom client struct using shared pool and single websocket instance. Each web client has one client connection on server
	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	// note: this is passing the client object to pool.register channel
	pool.Register <- client

	// start waiting to receive messages from web client
	client.Read()
}

func setupRoutes() {
	// create one pool for entire server
	pool := websocket.NewPool()

	// start the pool concurrently
	go pool.Start()

	// execute function when /ws route is hit
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

// start the server
func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()

	// serve to localhost
	http.ListenAndServe(":8081", nil)
}
