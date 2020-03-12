package observations_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

var coll *firestore.CollectionRef

// TestMain prepares the test suite for integration tests.
// It opens a connection to the database runs the tests and
// cleans up the connection.
//
// TODO: run firestore locally for integration tests.
func TestMain(m *testing.M) {
	var client *firestore.Client
	ctx := context.Background()

	var project = flag.String("project", "linked-data-land", "the google project the firestore instance is in.")

	flag.Parse()

	// Setup
	if !testing.Short() {
		client, err := firestore.NewClient(ctx, *project)
		if err != nil {
			fmt.Printf("connecting to firestore client: %v", err)
			os.Exit(1)
		}
		coll = client.Collection("observations_test")
	}

	// Run
	result := m.Run()

	// Teardown
	if !testing.Short() {
		if client != nil {
			client.Close()
		}
	}
	os.Exit(result)
}