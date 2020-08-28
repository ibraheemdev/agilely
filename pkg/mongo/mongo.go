package mongo

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client ...
type Client struct {
	client *mongo.Client
}

// GetClient ...
func GetClient(uri string) (engine.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return &Client{client}, err
}

// Connect ...
func (c *Client) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

// Ping ...
func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.client.Ping(ctx, rp)
}

// Disconnect ...
func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// Database ...
func (c *Client) Database(name string, opts ...*options.DatabaseOptions) engine.Database {
	return &Database{c.client.Database(name, opts...)}
}

// Database ...
type Database struct {
	database *mongo.Database
}

// Collection ...
func (m *Database) Collection(name string) engine.Collection {
	return &Collection{m.database.Collection(name)}
}

// Aggregate ...
func (m *Database) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.database.Aggregate(ctx, pipeline, opts...)
}

// Collection ...
type Collection struct {
	collection *mongo.Collection
}

// Name ...
func (c *Collection) Name() string {
	return c.collection.Name()
}

// Database ...
func (c *Collection) Database() engine.Database {
	return &Database{c.collection.Database()}
}

// Aggregate ...
func (c *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return c.collection.Aggregate(ctx, pipeline, opts...)
}

// Find ...
func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.collection.Find(ctx, filter, opts...)
}

// FindOne ...
func (c *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.collection.FindOne(ctx, filter, opts...)
}

// FindOneAndDelete ...
func (c *Collection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return c.collection.FindOneAndDelete(ctx, filter, opts...)
}

// FindOneAndUpdate ...
func (c *Collection) FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

// InsertMany ...
func (c *Collection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.collection.InsertMany(ctx, documents, opts...)
}

// InsertOne ...
func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document, opts...)
}

// UpdateOne ...
func (c *Collection) UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany ...
func (c *Collection) UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateMany(ctx, filter, update, opts...)
}

// DeleteMany ...
func (c *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.collection.DeleteMany(ctx, filter, opts...)
}

// DeleteOne ...
func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.collection.DeleteOne(ctx, filter, opts...)
}
