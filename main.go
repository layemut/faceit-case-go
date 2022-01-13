package main

import (
	"github.com/layemut/faceit-case-go/app"
)

func main() {
	application := &app.App{}
	application.Initialize()
	application.StartNotificationService()
	application.Run()
}
