package websocket

import (
	"fmt"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *ClientMessage
	Typing     chan *ClientMessage
	Login      chan *Client
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *ClientMessage),
		Typing:     make(chan *ClientMessage),
		Login:      make(chan *Client),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			fmt.Printf("client registered %s", client.ClientID)
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			client.Conn.WriteJSON(Client{ClientID: client.ClientID})
			break
		case client := <-pool.Unregister:
			_, ok := pool.Clients[client]
			fmt.Print("DELETING")
			if ok {
				delete(pool.Clients, client)
			}
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				fmt.Println(client.UserName, "client")
				sender := Client{ClientID: message.Client.ClientID, UserName: message.Client.UserName}
				if err := client.Conn.WriteJSON(ClientResponse{Client: &sender, Message: message.Message}); err != nil {
					fmt.Println(err)
					client.Conn.Close()
					delete(pool.Clients, client)
					break
				}
			}
			break
		case typing := <-pool.Typing:
			fmt.Println(typing.Typing)
			for client, _ := range pool.Clients {
				if typing.Client.ClientID != client.ClientID {
					if err := client.Conn.WriteJSON(ClientResponse{Client: typing.Client, Typing: typing.Typing}); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client)
						break
					}
				}
			}
			break
		case login := <-pool.Login:

			fmt.Printf("%s", login.UserName)
			for client, _ := range pool.Clients {
				if client.ClientID == login.ClientID {
					client.UserName = login.UserName
					fmt.Printf("\nFound Client ----> %s", login.UserName)
					if err := client.Conn.WriteJSON(Client{ClientID: login.ClientID, UserName: login.UserName}); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client)
						break
					}
				}
				if client.ClientID != login.ClientID {
					fmt.Println(login.UserName, login.ClientID, "<---- user \n")
					message := &ClientResponse{Client: &Client{UserName: login.UserName}, Message: "has joined..."}
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						client.Conn.Close()
						delete(pool.Clients, client)
						break
					}
				}
			}
			break
		}
	}
}
