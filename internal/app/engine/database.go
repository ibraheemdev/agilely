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
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	Collection(name string) Collection
}

// Collection ...
type Collection interface {
	Name() string
	Database() Database
	Drop(context.Context) error

	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)

	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult

	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult

	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
