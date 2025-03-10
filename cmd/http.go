package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/brionac626/taskManager/internal/repository"
	taskmanager "github.com/brionac626/taskManager/internal/taskManager"

	"github.com/spf13/cobra"
)

var port string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a task manager server",
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Starting server...")

		repo := repository.NewRepository()
		router := taskmanager.NewRouter(repo)
		go func() {
			if err := router.Start(":" + port); err != nil {
				log.Println("Error starting server", err)
				os.Exit(2)
			}
		}()

		// graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Println("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := router.Shutdown(ctx); err != nil {
			log.Println("shutdown server failed", err)
			return
		}

		log.Println("shutdown server successful")
	},
}
