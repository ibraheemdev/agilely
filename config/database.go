package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DatabaseClient : A pointer to the database client
	DatabaseClient *mongo.Database
)

// ConnectToDatabase :
func ConnectToDatabase() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", Config.Database.Host, Config.Database.Port)))
	if err != nil {
		log.Fatal(err)
	}
	DatabaseClient = client.Database(Config.Database.Name)
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
