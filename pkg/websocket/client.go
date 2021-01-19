package websocket

import (
	"fmt"
	"log"
)

// Client struct method: infinitely loops to wait to receive messages
func (c *Client) Read() {

	// this function will finally run when this Read() function is about to finish (when this object is about to be deleted)
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// blocks here until message arrives, then reads message and returns messageType and data
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// create HttpMessage object out of messageType and data
		// note: add sender's id to message
		body := ChatMessage{
			Username: c.Username,
			Message: string(p),
		}
		httpMessage := HttpMessage{
			Body: body,
		}
		broadcastChannelValue := BroadcastChannelValue{
			HttpMessage: httpMessage,
			CurrentClient: *c,
		}
		c.Pool.Broadcast <- broadcastChannelValue
		fmt.Printf("HttpMessage Received: %+v\n", httpMessage)
	}
}
