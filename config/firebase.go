package config

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Firebase struct {
	StorageClient *storage.Client
}

var firebaseInstance *Firebase

func InitFirebase() (*Firebase, error) {
	if firebaseInstance != nil {
		return firebaseInstance, nil
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return nil, fmt.Errorf("error initializing Google Cloud Storage client: %w", err)
	}

	firebaseInstance = &Firebase{StorageClient: client}
	return firebaseInstance, nil
}
