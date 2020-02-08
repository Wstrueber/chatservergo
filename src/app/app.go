package app

import (
	"chatservergo/src/app/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// RequestHandlerFunction ...
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

// App ...
type App struct {
	Router *mux.Router
}

// Init initializes Routers
func (a *App) Init() {
	a.Router = mux.NewRouter()
	a.SetRouters()
}

// SetRouters sets up routers
func (a *App) SetRouters() {
	a.Get("/ws", a.handleRequest(controllers.WebSocket))
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
