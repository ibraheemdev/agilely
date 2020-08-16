package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ibraheemdev/poller/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client : A pointer to the database client
	Client *mongo.Database
)

// Connect :
func Connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Config.Database.Host, config.Config.Database.Port)))
	if err != nil {
		log.Fatal(err)
	}
	Client = client.Database(config.Config.Database.Name)
	return client
}

// Disconnect :
func Disconnect(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
