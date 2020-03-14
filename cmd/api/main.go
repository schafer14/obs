package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/schafer14/observations/cmd/api/internal/handlers"
	"github.com/schafer14/observations/internal/platform/database"
)

var build = "develop"

func main() {
	if err := run(); err != nil {
		fmt.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// =============================================== //
	// Read Configuration
	// =============================================== //
	var cfg struct {
		APIHost          string `conf:"default:0.0.0.0:3000"`
		FirestoreProject string `conf:"default:linked-data-land"`
		WithDocs         bool   `conf:"default:true"`
	}

	if err := conf.Parse(os.Args[1:], "OBS", &cfg); err != nil {
		fmt.Println(err)
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("OBS", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =============================================== //
	// Report Build Parameters
	// =============================================== //
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =============================================== //
	// Configure database
	// =============================================== //
	log.Println("main : Started : Initializing arango database support")

	db, err := database.Open(ctx, cfg.FirestoreProject)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}

	// =============================================== //
	// Starting API
	// =============================================== //
	log.Println("main : Started : Initializing API support")

	router := handlers.API(db)

	http.ListenAndServe(cfg.APIHost, router)

	return nil
}
