package main

import (
	"log"
	"os"

	"github.com/brionac626/taskManager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println("Execute command error", err)
		os.Exit(5)
	}
}
