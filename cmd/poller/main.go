package main

import (
	"log"
	"os"

	"github.com/ibraheemdev/poller/config"
	"github.com/ibraheemdev/poller/internal/router"
	"github.com/ibraheemdev/poller/pkg/database"
)

func main() {
	env := os.Getenv("POLLER_ENV")
	if env != "testing" && env != "development" && env != "production" {
		log.Fatal("must set POLLER_ENV to a valid environment")
	}
	log.Printf("starting application in %s environment", env)

	config.Init()
	client := database.Connect()
	defer database.Disconnect(client)
	router.ListenAndServe()
}
