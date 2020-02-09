package app

import (
	socket "chatservergo/src/app/websocket"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// RequestHandlerFunction ...
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

// App ...
type App struct {
	Router *mux.Router
}

// Init initializes Routers
func (a *App) Init(pool *socket.Pool) {
	a.Router = mux.NewRouter()
	a.SetupRoutes(pool)
}

func serveWs(pool *socket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &socket.Client{
		Conn:     conn,
		Pool:     pool,
		ClientID: uuid.New(),
	}
	fmt.Println(client.ClientID)
	pool.Register <- client
	client.Read()
}

// SetupRoutes sets up routers
func (a *App) SetupRoutes(pool *socket.Pool) {
	a.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

// Run runs the server
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

// Get handles get methods
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods(http.MethodGet)
}
