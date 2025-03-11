package main

import (
	"log"
	"os"

	"github.com/brionac626/taskManager/cmd"
)

// @title Task Manager API
// @version 0.1.2
// @description This is a sample task manager server.

// @contact.name API support
// @contact.url https://github.com/brionac626/taskManager
// @contact.email theone1632@gmail.com

// @license.name MIT
// @license.url https://github.com/brionac626/taskManager/blob/main/LICENSE

// @host localhost:8080
// @BasePath

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println("Execute command error", err)
		os.Exit(5)
	}
}
