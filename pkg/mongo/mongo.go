package mongo

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client ...
type Client struct {
	*mongo.Client
}

// GetClient ...
func GetClient(uri string) (engine.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return &Client{client}, err
}

// Database ...
func (c *Client) Database(name string, opts ...*options.DatabaseOptions) engine.Database {
	return &Database{c.Client.Database(name, opts...)}
}

// Database ...
type Database struct {
	*mongo.Database
}

// Collection ...
func (m *Database) Collection(name string) engine.Collection {
	return &Collection{m.Database.Collection(name)}
}

// Collection ...
type Collection struct {
	*mongo.Collection
}

// Database ...
func (c *Collection) Database() engine.Database {
	return &Database{c.Collection.Database()}
}
