package websocket

import (
	"fmt"
	// "strconv"
)

// create a new pool instance
func NewPool() *Pool {
	return &Pool{
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan BroadcastChannelValue),
		NextId: 1,
	}
}

// Pool struct method: infinitely loop to wait for and execute the pool command
func (pool *Pool) Start() {
	for {
		select {

		// if register channel has received message
		case newClient := <-pool.Register :

			// assign id to client and increment
			newClient.Id = pool.NextId
			pool.NextId++

			// generate default username to new client
			newClient.Username = fmt.Sprintf("User %d", newClient.Id)

			// place newClient into pool's map of clients
			pool.Clients[newClient] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			// when client registers, send a message to new client only providing their Id
			// note: sending message as a "system message denoted by fromclientid == 0"
			newClient.Conn.WriteJSON(HttpMessage{ Body: ChatMessage{ Message: fmt.Sprintf("Welcome %v", newClient.Username) } })
			for client, _ := range pool.Clients { // for each client in the pool
				if newClient.Id != client.Id {
					client.Conn.WriteJSON(HttpMessage{ Body: ChatMessage{ Message: fmt.Sprintf("%v Joined...", newClient.Username) } }) // send "new user" to each client as "system message"
				}
			}
			break

		// if unregister channel has received message
		case targetClient := <-pool.Unregister :
			delete(pool.Clients, targetClient) // remove the unregistered client from the pool of clients
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(HttpMessage{ Body: ChatMessage{ Message: fmt.Sprintf("%v Disconnected...", targetClient.Username) } }) // send "user disconnected" to each client as "system message"
			}
			break

		// if broadcast channel has received message
		// message should contain the sender's id
		case broadcastChannelValue := <-pool.Broadcast :
			fmt.Println("Sending message to all clients in Pool")
			currentClient := broadcastChannelValue.CurrentClient
			httpMessage := broadcastChannelValue.HttpMessage
			for client, _ := range pool.Clients {
				httpMessage.Body.IsOwner = client.Id == currentClient.Id
				if err := client.Conn.WriteJSON(httpMessage); err != nil { // send message to each client
					fmt.Println(err)
					return
				}
			}
		}
	}
}
