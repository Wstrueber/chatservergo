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
	fmt.Printf("Server has started on port %s", port)
	app.Run(":" + port)
}
