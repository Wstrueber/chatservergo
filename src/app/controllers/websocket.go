package controllers

import (
	"chatservergo/src/app/constants"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ClientMessage ...
type ClientMessage struct {
	Message string `"json":"message"`
}

// JSONResponse ...
type JSONResponse struct {
	AppVersionNumber string      `json:"appVersionNumber"`
	Data             interface{} `json:"data,omitempty"`
}

// ClientResponse ...
type ClientResponse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

var conn *websocket.Conn

func writeJSON(clientResponse *ClientResponse) error {
	cr := JSONResponse{AppVersionNumber: constants.AppVersionNumber}
	logMsg := fmt.Sprintf("websocket.writeJson: writing to client: app version number: %q", cr.AppVersionNumber)
	if clientResponse != nil {
		cr.Data = &clientResponse
		logMsg = fmt.Sprintf("websocket.writeJson: writing to client: app version number: %q, message: %q", cr.AppVersionNumber, clientResponse.Message)
	}
	err := conn.WriteJSON(cr)
	if err != nil {
		err = fmt.Errorf("websocket.writeJSON: Failed to send JSON %w", err)
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(logMsg)
	return nil
}

func reader() error {
	clientMessage := new(ClientMessage)
	for {
		err := conn.ReadJSON(&clientMessage)
		if err != nil {
			err = fmt.Errorf("websocket.reader: Failed to read client message %w", err)
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(fmt.Sprintf("websocket.reader: Client message: %q", clientMessage.Message))
		if clientMessage.Message == constants.RequestVersionNumber {
			writeJSON(nil)
		} else {
			err := writeJSON(&ClientResponse{
				Message: "Server has received the response",
				Time:    time.Now().UTC(),
			})
			if err != nil {
				err = fmt.Errorf("websocket.reader: Failed to send message to client: %w", err)
				fmt.Println(err.Error())
				return err
			}
		}
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
		return true
		// return r.Header.Get("Origin") == "http://localhost:3000"
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("websocket.WebSocket: Failed to upgrade the HTTP server connection to the WebSocket protocol, %s", err.Error())
		writeJSON(&ClientResponse{Message: err.Error()})
		return
	}
	conn = ws

	fmt.Println("websocket.WebSocket: Client Successfully Connected...")
	ticker := createTicker(10 * time.Second)
	ticker.pushVersionNumber()
	err = reader()
	if err != nil {
		fmt.Printf("websocket.WebSocket: Failed to read client message %s", err.Error())
		writeJSON(&ClientResponse{Message: err.Error()})
		return
	}
}
