package database

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

func Open(ctx context.Context, projectId string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		return nil, errors.Wrap(err, "connecting to firestore")
	}

	return client, nil
}
