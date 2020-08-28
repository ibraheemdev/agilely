package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ibraheemdev/agilely/config"
	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/pkg/mongo"
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

	client, err := mongo.GetClient(fmt.Sprintf("mongodb://%s:%d", e.Config.Database.Host, e.Config.Database.Port))
	if err != nil {
		return err
	}
	e.Core.Database = client.Database(e.Config.Database.Name)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	http.ListenAndServe(fmt.Sprintf("%s:%d", e.Config.Server.Host, e.Config.Server.Port), e.Core.Router)

	return nil
}
