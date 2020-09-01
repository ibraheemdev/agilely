package mongo

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// context timeout for database queries
	queryTimeout = 5 * time.Second
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
	collection *mongo.Collection
}

// Name ...
func (c *Collection) Name() string {
	return c.collection.Name()
}

// Drop ...
func (c *Collection) Drop(ctx context.Context) error {
	return c.collection.Drop(ctx)
}

// Aggregate ...
func (c *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.Aggregate(ctx, pipeline, opts...)
}

// Find ...
func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.Find(ctx, filter, opts...)
}

// FindOne ...
func (c *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOne(ctx, filter, opts...)
}

// FindOneAndDelete ...
func (c *Collection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOneAndDelete(ctx, filter, opts...)
}

// FindOneAndUpdate ...
func (c *Collection) FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

// InsertMany ...
func (c *Collection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.InsertMany(ctx, documents, opts...)
}

// InsertOne ...
func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.InsertOne(ctx, document, opts...)
}

// UpdateOne ...
func (c *Collection) UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany ...
func (c *Collection) UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.UpdateMany(ctx, filter, update, opts...)
}

// DeleteMany ...
func (c *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.DeleteMany(ctx, filter, opts...)
}

// DeleteOne ...
func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.DeleteOne(ctx, filter, opts...)
}

// Database ...
func (c *Collection) Database() engine.Database {
	return &Database{c.collection.Database()}
}
