package websocket

import (
	"chatservergo/src/app/constants"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn `json:"-"`
	Pool     *Pool           `json:"-"`
	UserName string          `json:"userName"`
	ClientID uuid.UUID       `json:"clientId"`
}

// JSONResponse ...
type JSONResponse struct {
	AppVersionNumber string `json:"appVersionNumber,omitempty"`
}

// ClientMessage ...
type ClientMessage struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Client  *Client `json:"client"`
	Typing  bool    `json:"typing,omitempty"`
}

// ClientResponse ...
type ClientResponse struct {
	*Client
	Message string     `json:"message,omitempty"`
	Time    *time.Time `json:"serverTime,omitempty"`
	Typing  bool       `json:"typing,omitempty"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var clientMessage *ClientMessage
		c.Conn.ReadJSON(&clientMessage)
		c.handleClientRequest(clientMessage)
	}
}

func (c *Client) handleClientRequest(clientMessage *ClientMessage) {
	switch clientMessage.Action {
	case constants.RequestResponse:
		fmt.Printf("client message %q", clientMessage.Client.UserName)
		c.Pool.Broadcast <- clientMessage
		return
	case constants.RequestVersionNumber:
		c.Conn.WriteJSON(JSONResponse{AppVersionNumber: constants.AppVersionNumber})
		return
	case constants.UserTyping:
		fmt.Printf("Client ---> %s %s\n", clientMessage.Client.ClientID, clientMessage.Client.UserName, clientMessage.Typing)
		c.Pool.Typing <- clientMessage
		return
	case constants.RequestLogin:
		c.Pool.Login <- clientMessage.Client
		return
	}
	return
}
