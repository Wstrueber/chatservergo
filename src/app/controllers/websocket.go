package controllers

import (
	"chatservergo/src/app/constants"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ClientMessage ...
type ClientMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

// JSONResponse ...
type JSONResponse struct {
	AppVersionNumber string      `json:"appVersionNumber"`
	Response         interface{} `json:"response,omitempty"`
}
type Client struct {
	UserName string    `json:"userName"`
	ClientID uuid.UUID `json:"clientId"`
}

// ClientResponse ...
type ClientResponse struct {
	*Client
	JSONResponse
	Data string     `json:"data,omitempty"`
	Time *time.Time `json:"serverTime,omitempty"`
}

var conn *websocket.Conn

var client *Client

func writeJSON(clientResponse *ClientResponse) error {

	logMsg := fmt.Sprintf("websocket.writeJson: writing to client: app version number: %q", constants.AppVersionNumber)
	if clientResponse != nil {
		logMsg = fmt.Sprintf("websocket.writeJson: writing to client: app version number: %q, message: %q", constants.AppVersionNumber, clientResponse.Data)
	}

	err := conn.WriteJSON(clientResponse)
	if err != nil {
		err = fmt.Errorf("websocket.writeJSON: Failed to send JSON %w", err)
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(logMsg)
	return nil
}

func handleClientRequest(clientMessage *ClientMessage, clientID uuid.UUID) {
	switch clientMessage.Action {
	case constants.RequestResponse:
		fmt.Printf("made it here --- %s", client.UserName)
		err := writeJSON(&ClientResponse{
			Client: client,
			JSONResponse: JSONResponse{
				AppVersionNumber: constants.AppVersionNumber,
				Response:         map[string]interface{}{"message": clientMessage.Message},
			},
		})
		if err != nil {
			return
		}
		return
	case constants.RequestVersionNumber:
		err := writeJSON(&ClientResponse{
			JSONResponse: JSONResponse{
				AppVersionNumber: constants.AppVersionNumber,
			},
		})
		if err != nil {
			return
		}
		return
	case constants.UserTyping:
		err := writeJSON(&ClientResponse{
			Client: client,
			JSONResponse: JSONResponse{
				AppVersionNumber: constants.AppVersionNumber,
				Response:         map[string]interface{}{"typing": true},
			},
		})
		if err != nil {
			return
		}
		return
	case constants.RequestLogin:
		client = &Client{UserName: clientMessage.Message, ClientID: clientID}
		fmt.Printf("%s", client.UserName)
		err := writeJSON(&ClientResponse{
			Client: client,
			JSONResponse: JSONResponse{
				AppVersionNumber: constants.AppVersionNumber,
				Response:         map[string]interface{}{"userName": clientMessage.Message},
			},
		})
		if err != nil {
			return
		}
		return
	}
	return
}

func reader(clientID uuid.UUID) error {
	clientMessage := new(ClientMessage)

	for {
		err := conn.ReadJSON(&clientMessage)
		if err != nil {
			err = fmt.Errorf("websocket.reader: Failed to read client message %w", err)
			fmt.Println(err.Error())
			return err
		}
		handleClientRequest(clientMessage, clientID)
		fmt.Println(fmt.Sprintf("websocket.reader: Client message: %q", clientMessage.Message))
	}
}

// Ticker ...
type Ticker struct {
	ticker   time.Ticker
	duration time.Duration
}

func createTicker(duration time.Duration) *Ticker {
	return &Ticker{ticker: *time.NewTicker(duration), duration: duration}
}

func (t *Ticker) pushVersionNumber() {
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-t.ticker.C:
				fmt.Println("made it here writing json")
				writeJSON(nil)
			case <-quit:
				fmt.Println("made it here stopping ticker")
				t.ticker.Stop()
				return
			}
		}
	}()
}

// WebSocket controller
func WebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// return true
		return r.Header.Get("Origin") == "http://localhost:3000"
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("websocket.WebSocket: Failed to upgrade the HTTP server connection to the WebSocket protocol, %s", err.Error())
		writeJSON(&ClientResponse{Data: err.Error()})
		return
	}
	conn = ws

	fmt.Println("websocket.WebSocket: Client Successfully Connected...")
	clientID := uuid.New()
	err = reader(clientID)
	if err != nil {
		fmt.Printf("websocket.WebSocket: Failed to read client message %s", err.Error())
		writeJSON(&ClientResponse{Data: err.Error()})
		return
	}
}
