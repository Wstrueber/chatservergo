package main

import (
	"chatservergo/src/app"
	"fmt"
)

func main() {
	app := &app.App{}
	app.Init()
	fmt.Println("Server has started on port 8080...")
	app.Run(":8080")
}
