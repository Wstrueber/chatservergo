package websocket

import (
	"chatservergo/src/app/utils"
	"fmt"

	"github.com/google/uuid"
)

// Pool pool type
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID]*Client
	Broadcast  chan *ClientMessage
	Typing     chan *ClientMessage
	Login      chan *Client
}

// NewPool creats the pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan *ClientMessage),
		Typing:     make(chan *ClientMessage),
		Login:      make(chan *Client),
	}
}

// Start starts the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client.ClientID] = client
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			client.Conn.WriteJSON(Client{ClientID: client.ClientID})
			break
		case client := <-pool.Unregister:
			_, ok := pool.Clients[client.ClientID]
			if ok {
				delete(pool.Clients, client.ClientID)
			}
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			sender := pool.Clients[message.Client.ClientID]

			clients := make(Clients, 0)
			for _, value := range pool.Clients {
				if value.ClientID != uuid.Nil {
					clients.append(value)
				}
			}
			for _, client := range pool.Clients {
				if err := client.Conn.WriteJSON(ClientResponse{Client: sender, Message: message.Message, Clients: clients}); err != nil {
					fmt.Println(err)
					client.Conn.Close()
					delete(pool.Clients, client.ClientID)
					break
				}
			}
			break
		case typing := <-pool.Typing:
			fmt.Println(typing.Typing)
			for _, client := range pool.Clients {
				clients := make(Clients, 0)
				for _, value := range pool.Clients {
					if value.ClientID != uuid.Nil {
						clients.append(value)
					}
				}
				if typing.Client.ClientID != client.ClientID {
					if err := client.Conn.WriteJSON(ClientResponse{Client: pool.Clients[typing.Client.ClientID], Clients: clients, Typing: typing.Typing}); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client.ClientID)
						break
					}
				}
			}
			break
		case login := <-pool.Login:
			for _, client := range pool.Clients {
				clients := make(Clients, 0)
				for _, value := range pool.Clients {
					if value.ClientID != uuid.Nil {
						clients.append(value)
					}
				}
				if client.ClientID == login.ClientID {
					client.UserName = login.UserName
					client.Color = utils.GetRandomColorInHex()
					if err := client.Conn.WriteJSON(Client{ClientID: login.ClientID, Color: client.Color, UserName: login.UserName}); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client.ClientID)
						break
					}
				}
				if client.ClientID != login.ClientID {
					color := pool.Clients[login.ClientID].Color
					message := &ClientResponse{Client: &Client{UserName: login.UserName, Color: color}, Clients: clients, Message: "has joined..."}
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client.ClientID)
						break
					}
				}
			}
			break
		}
	}
}
