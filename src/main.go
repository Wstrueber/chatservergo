package main

import (
	"chatservergo/src/app"
	socket "chatservergo/src/app/websocket"
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pool := socket.NewPool()
	go pool.Start()

	app := &app.App{}
	app.Init(pool)
	fmt.Printf("Server has started on port :%s", port)
	app.Run(":" + port)
}
