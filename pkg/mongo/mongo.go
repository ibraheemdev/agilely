package mongo

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Find ...
func (c *Collection) Find(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOptions) (engine.Cursor, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.Find(ctx, bson.M(filter), opts...)
}

// FindOne ...
func (c *Collection) FindOne(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneOptions) engine.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOne(ctx, bson.M(filter), opts...)
}

// FindOneAndDelete ...
func (c *Collection) FindOneAndDelete(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneAndDeleteOptions) engine.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOneAndDelete(ctx, bson.M(filter), opts...)
}

// FindOneAndUpdate ...
func (c *Collection) FindOneAndUpdate(ctx context.Context, filter, update map[string]interface{}, opts ...*options.FindOneAndUpdateOptions) engine.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	return c.collection.FindOneAndUpdate(ctx, bson.M(filter), bson.D{{"$set", bson.M(update)}}, opts...)
}

// InsertMany ...
func (c *Collection) InsertMany(ctx context.Context, documents []map[string]interface{}, opts ...*options.InsertManyOptions) (engine.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	docs := make([]interface{}, len(documents))
	for i, v := range documents {
		docs[i] = bson.M(v)
	}
	res, err := c.collection.InsertMany(ctx, docs, opts...)
	return &insertManyResult{res}, err
}

type insertManyResult struct {
	*mongo.InsertManyResult
}

func (i *insertManyResult) InsertedIDs() []string {
	var ids []string
	for _, i := range i.InsertManyResult.InsertedIDs {
		ids = append(ids, i.(primitive.ObjectID).Hex())
	}
	return ids
}

// InsertOne ...
func (c *Collection) InsertOne(ctx context.Context, document map[string]interface{}, opts ...*options.InsertOneOptions) (engine.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	res, err := c.collection.InsertOne(ctx, bson.M(document), opts...)
	return &insertOneResult{res}, err
}

type insertOneResult struct {
	*mongo.InsertOneResult
}

func (i *insertOneResult) InsertedID() string {
	return i.InsertOneResult.InsertedID.(primitive.ObjectID).Hex()
}

// UpdateOne ...
func (c *Collection) UpdateOne(ctx context.Context, filter, update map[string]interface{}, opts ...*options.UpdateOptions) (engine.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	res, err := c.collection.UpdateOne(ctx, bson.M(filter), bson.D{{"$set", bson.M(update)}}, opts...)
	return &updateResult{res}, err
}

// UpdateMany ...
func (c *Collection) UpdateMany(ctx context.Context, filter, update map[string]interface{}, opts ...*options.UpdateOptions) (engine.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	res, err := c.collection.UpdateMany(ctx, bson.M(filter), bson.D{{"$set", bson.M(update)}}, opts...)
	return &updateResult{res}, err
}

type updateResult struct {
	*mongo.UpdateResult
}

func (u *updateResult) MatchedCount() int64 {
	return u.UpdateResult.MatchedCount
}

func (u *updateResult) ModifiedCount() int64 {
	return u.UpdateResult.ModifiedCount
}

func (u *updateResult) UpsertedCount() int64 {
	return u.UpdateResult.UpsertedCount
}

func (u *updateResult) UpsertedID() string {
	return u.UpdateResult.UpsertedID.(primitive.ObjectID).Hex()
}

// DeleteMany ...
func (c *Collection) DeleteMany(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (engine.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	res, err := c.collection.DeleteMany(ctx, bson.M(filter), opts...)
	return &deleteResult{res}, err
}

// DeleteOne ...
func (c *Collection) DeleteOne(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (engine.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	res, err := c.collection.DeleteOne(ctx, bson.M(filter), opts...)
	return &deleteResult{res}, err
}

type deleteResult struct {
	*mongo.DeleteResult
}

func (d *deleteResult) DeletedCount() int64 {
	return d.DeleteResult.DeletedCount
}

// Database ...
func (c *Collection) Database() engine.Database {
	return &Database{c.collection.Database()}
}
