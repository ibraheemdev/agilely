package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ibraheemdev/agilely/config"
	"github.com/ibraheemdev/agilely/internal/app/engine"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	env := os.Getenv("AGILELY_ENV")
	if env != "testing" && env != "development" && env != "production" {
		return fmt.Errorf("must set AGILELY_ENV to a valid environment")
	}
	log.Printf("starting application in %s environment", env)

	e := engine.New()
	config.SetCore(e)

	if err := config.SetConfig(e); err != nil {
		return err
	}

	if err := config.Routes(e); err != nil {
		return err
	}

	client := config.ConnectToDatabase(e)
	defer config.DisconnectFromDatabase(client)

	http.ListenAndServe(fmt.Sprintf("%s:%d", e.Config.Server.Host, e.Config.Server.Port), e.Core.Router)

	return nil
}
