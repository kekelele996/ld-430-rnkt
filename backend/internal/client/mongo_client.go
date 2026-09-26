package client

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient wraps the mongo-driver client and database.
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoClient connects to MongoDB and returns a ready wrapper.
func NewMongoClient(ctx context.Context, uri, database string) (*MongoClient, error) {
	opts := options.Client().ApplyURI(uri).SetServerSelectionTimeout(5 * time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}
	db := client.Database(database)
	return &MongoClient{Client: client, DB: db}, nil
}

// Close disconnects the underlying client.
func (m *MongoClient) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// Collection returns a named collection.
func (m *MongoClient) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
