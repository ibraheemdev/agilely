package main

import (
	"log"
	"os"

	"github.com/ibraheemdev/agilely/config"
	"github.com/ibraheemdev/agilely/internal/app/router"
	_ "github.com/ibraheemdev/agilely/internal/app/users"
	"github.com/ibraheemdev/agilely/pkg/database"
	"github.com/julienschmidt/httprouter"
)

func main() {
	env := os.Getenv("AGILELY_ENV")
	if env != "testing" && env != "development" && env != "production" {
		log.Fatal("must set AGILELY_ENV to a valid environment")
	}
	log.Printf("starting application in %s environment", env)

	config.Init()
	client := database.Connect()
	defer database.Disconnect(client)
	r := httprouter.New()
	config.SetupAuthboss(r)
	router.ListenAndServe(r)
}
