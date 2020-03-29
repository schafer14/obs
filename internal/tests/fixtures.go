package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/schafer14/obs/internal/platform/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container tracks information about a docker container started for tests.
type Container struct {
	ID   string
	Host string // IP:Port
}

func DatabaseTest(t *testing.T, c *Container) *mongo.Database {
	t.Helper()

	ctx := context.Background()

	maxAttempts := 20
	var db *mongo.Database
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		dbTry, err := database.Open(ctx, c.Host, "observations")

		if err == nil {
			db = dbTry
			break
		}
		if err != nil && attempts == maxAttempts {
			t.Fatalf("opening database connection: %v", err)
		}
		time.Sleep(time.Second)
	}

	t.Log("waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = database.Check(ctx, db.Client())
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		dumpContainerLogs(t, c)
		TeardownDatabase(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}

	return db
}

// StartContainer runs a postgres container to execute commands.
func SetupDatabase(t *testing.T) *Container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-P", "-d", "mongo:4-bionic")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start container: %v", err)
	}

	id := out.String()[:12]
	t.Log("DB ContainerID:", id)

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container %s: %v", id, err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP5432 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"27017/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP5432[0]

	c := Container{
		ID:   id,
		Host: "mongodb://" + network.HostIP + ":" + network.HostPort,
	}

	t.Log("DB Host:", c.Host)

	return &c
}

// StopContainer stops and removes the specified container.
func TeardownDatabase(t *testing.T, c *Container) {
	t.Helper()

	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}
	t.Log("Stopped:", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}
	t.Log("Removed:", c.ID)
}

// DumpContainerLogs runs "docker logs" against the container and send it to t.Log
func dumpContainerLogs(t *testing.T, c *Container) {
	t.Helper()

	out, err := exec.Command("docker", "logs", c.ID).CombinedOutput()
	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}
	t.Logf("Logs for %s\n%s:", c.ID, out)
}
