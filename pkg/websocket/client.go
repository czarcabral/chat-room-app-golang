package websocket

import (
	"fmt"
	"log"
	// "sync"
	"github.com/gorilla/websocket"
)

// Each client is a websocket connection with an Id and access to shared pool
type Client struct {
	ID 		int
	Conn 	*websocket.Conn
	Pool 	*Pool
}

// Each request/response exists in this form
type Message struct {
	Type int `json:"type"`
	Body string `json:"body"`
	FromClientId int `json:"fromClientId"` // if fromClientId == 0, it means that the message is a system message like Register/Unregister
	NewClientId int `json:"newClientId"` // this will hold the new client's id to be sent once they register, but will be 0 for every other message
}

// Client struct method: infinitely loops to wait to receive messages
func (c *Client) Read() {

	// this function will finally run when this Read() function is about to finish (when this object is about to be deleted)
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// blocks here until message arrives, then reads message and returns messageType and data
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// create Message object out of messageType and data
		message := Message{Type: messageType, Body: string(p), FromClientId: c.ID}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
