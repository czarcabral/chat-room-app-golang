package websocket

import (
	"fmt"
	// "strconv"
)

// can register/unregister a client, contains the map of clients, and broadcast method to send a message to all clients
// note: channels are used so that while pool Start is running in parallel, it is waiting to receive something in any of its channels before it does something
type Pool struct {
	Register chan *Client
	Unregister chan *Client
	Clients map[*Client]bool
	Broadcast chan Message
	NextId int
}

// create a new pool instance
func NewPool() *Pool {
	return &Pool {
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan Message),
		NextId: 1,
	}
}

// Pool struct method: infinitely loop to wait for and execute the pool command
func (pool *Pool) Start() {
	for {
		select {

		// if Register channel has received message
		case newClient := <-pool.Register :
			newClient.ID = pool.NextId
			pool.NextId++
			pool.Clients[newClient] = true // place newClient into pool's map of clients
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			newClient.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("Welcome. You are User %d.", newClient.ID), FromClientId: 0, NewClientId: newClient.ID})
			for client, _ := range pool.Clients { // for each client in the pool
				if newClient.ID != client.ID {
					client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("User %d Joined...", newClient.ID), FromClientId: 0}) // send "new user" to each client
				}
			}
			break

		// if Unregister channel has received message
		case targetClient := <-pool.Unregister :
			delete(pool.Clients, targetClient) // remove the unregistered client from the pool of clients
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("User %d Disconnected...", targetClient.ID), FromClientId: 0}) // send "user disconnected" to each client
			}
			break

		// if Broadcast channel has received message
		case message := <-pool.Broadcast :
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil { // send message to each client
					fmt.Println(err)
					return
				}
			}
		}
	}
}
