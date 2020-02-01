package main

import (
	"chatservergo/src/app"
)

func main() {
	app := &app.App{}
	app.Init()
	app.Run(":8080")
}
