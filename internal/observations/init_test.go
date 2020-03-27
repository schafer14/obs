package observations_test

import (
	"flag"
	"os"
	"testing"

	"github.com/schafer14/observations/internal/tests"
	"go.mongodb.org/mongo-driver/mongo"
)

var coll *mongo.Collection

// TestMain runs a database for this package.
func TestMain(m *testing.M) {
	t := &testing.T{}

	flag.Parse()
	var c *tests.Container
	if !testing.Short() {
		c = tests.SetupDatabase(t)
		db := tests.DatabaseTest(t, c)
		coll = db.Collection("observations")
	}

	result := m.Run()

	if !testing.Short() {
		tests.TeardownDatabase(t, c)
	}
	os.Exit(result)
}
