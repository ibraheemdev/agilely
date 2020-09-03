package engine

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// ErrNoDocuments ...
	ErrNoDocuments = mongo.ErrNoDocuments
)

// Client ...
type Client interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error

	Disconnect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) Database
}

// Database ...
type Database interface {
	Collection(name string) Collection
}

// Collection ...
type Collection interface {
	Name() string
	Database() Database
	Drop(context.Context) error

	Find(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOptions) (Cursor, error)
	FindOne(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneOptions) SingleResult

	FindOneAndDelete(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResult
	FindOneAndUpdate(ctx context.Context, filter map[string]interface{}, update map[string]interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResult

	InsertMany(ctx context.Context, documents []map[string]interface{}, opts ...*options.InsertManyOptions) (InsertManyResult, error)
	InsertOne(ctx context.Context, document map[string]interface{}, opts ...*options.InsertOneOptions) (InsertOneResult, error)

	UpdateMany(ctx context.Context, filter, update map[string]interface{}, opts ...*options.UpdateOptions) (UpdateResult, error)
	UpdateOne(ctx context.Context, filter, update map[string]interface{}, opts ...*options.UpdateOptions) (UpdateResult, error)

	DeleteMany(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (DeleteResult, error)
	DeleteOne(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (DeleteResult, error)
}

// SingleResult ...
type SingleResult interface {
	Err() error
	Decode(interface{}) error
}

// Cursor ...
type Cursor interface {
	Err() error
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	All(context.Context, interface{}) error
}

// InsertManyResult ...
type InsertManyResult interface {
	InsertedIDs() []string
}

// InsertOneResult ...
type InsertOneResult interface {
	InsertedID() string
}

// UpdateResult ...
type UpdateResult interface {
	MatchedCount() int64
	ModifiedCount() int64
	UpsertedCount() int64
	UpsertedID() string
}

// DeleteResult ...
type DeleteResult interface {
	DeletedCount() int64
}
