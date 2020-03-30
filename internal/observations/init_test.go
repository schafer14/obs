package observations_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/schafer14/obs/internal/tests"
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
		db, err := tests.DatabaseTest(t, c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		coll = db.Collection("observations")
	}

	result := m.Run()

	if !testing.Short() {
		tests.TeardownDatabase(t, c)
	}
	os.Exit(result)
}
