package main

import (
	"log"
	"os"

	"github.com/ibraheemdev/poller/config"
	"github.com/ibraheemdev/poller/internal/app/router"
	_ "github.com/ibraheemdev/poller/internal/app/users"
	"github.com/ibraheemdev/poller/pkg/database"
	"github.com/julienschmidt/httprouter"
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
	r := httprouter.New()
	config.SetupAuthboss(r)
	router.ListenAndServe(r)
}
