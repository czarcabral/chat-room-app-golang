package websocket

import (
	"github.com/gorilla/websocket"
)

// can register/unregister a client, contains the map of clients, and broadcast method to send a message to all clients
// note: channels are used so that while pool Start is running in parallel, it is waiting to receive something in any of its channels before it does something
type Pool struct {
	Register chan *Client
	Unregister chan *Client
	Clients map[*Client]bool
	Broadcast chan BroadcastChannelValue
	NextId int // the shared pool will assign the clients Ids when they are first registered
}

// Each client is a websocket connection with an Id and access to shared pool
type Client struct {
	Id int
	Username string
	Conn *websocket.Conn
	Pool *Pool
}

// Our custom HttpMessage (request/response)
type HttpMessage struct {
	Body ChatMessage `json:"body"`
}

// Our ChatMessage object
type ChatMessage struct {
	Username string `json:"username"`
	Message string `json:"message"`
	IsOwner bool `json:"isOwner"`
}

type BroadcastChannelValue struct {
	HttpMessage HttpMessage
	CurrentClient Client
}

// // ChatEventType enum for determining what event is happening in message
// type ChatEventType int
// const (
// 	Registering ChatEventType = iota
// 	Unregistering
// 	UsernameChange
// 	Default
// )
// func (ce ChatEventType) String() string {
// 	return [...]string{"Registering", "Unregistering", "UsernameChange", "Default"}[ce] // note: this syntax is creating array literal and then accessing element at index ce
// }
