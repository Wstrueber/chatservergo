package main

import (
	"chatservergo/src/app"
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	app := &app.App{}
	app.Init()
	fmt.Println("Server has started on port 8080...")
	app.Run(port)
}
