package database

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Open opens a connection to a mongo database
func Open(ctx context.Context, connectionString string, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	return client.Database(database), errors.Wrap(err, "connecting to database")
}

// Check makes sure the db connection is responding.
func Check(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())

	return errors.Wrap(err, "pinging database")
}
