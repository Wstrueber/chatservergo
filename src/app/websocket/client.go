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
	Color    string          `json:"color"`
}

type Clients []*Client

func (c *Clients) append(client *Client) {
	*c = append(*c, client)
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
	Clients []*Client  `json:"Clients"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var clientMessage *ClientMessage
		err := c.Conn.ReadJSON(&clientMessage)
		if err != nil {
			err = fmt.Errorf("client.Read: Error reading json, %s", err.Error())
			fmt.Println(err)
			return
		}

		if clientMessage != nil {
			switch clientMessage.Action {
			case constants.RequestResponse:
				c.Pool.Broadcast <- clientMessage
			case constants.RequestVersionNumber:
				c.Conn.WriteJSON(JSONResponse{AppVersionNumber: constants.AppVersionNumber})
			case constants.UserTyping:
				c.Pool.Typing <- clientMessage
			case constants.RequestLogin:
				c.Pool.Login <- clientMessage.Client
			}
		}

	}
}
