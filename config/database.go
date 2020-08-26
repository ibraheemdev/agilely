package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DatabaseClient : A pointer to the database client
	DatabaseClient *mongo.Database
)

// ConnectToDatabase :
func ConnectToDatabase(e *engine.Engine) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", e.Config.Database.Host, e.Config.Database.Port)))
	if err != nil {
		log.Fatal(err)
	}
	DatabaseClient = client.Database(e.Config.Database.Name)
	return client
}

// DisconnectFromDatabase :
func DisconnectFromDatabase(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
