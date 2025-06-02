package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"pelita/config"
	"pelita/entity"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func UploadFile(user_id uuid.UUID, ctx string, file *multipart.FileHeader, fileExt string) (string, error) {
	firebase, err := config.InitFirebase()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Firebase: %w", err)
	}
	bucket := firebase.StorageClient.Bucket(os.Getenv("FIREBASE_BUCKET_NAME"))
	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	id := uuid.New().String()
	objectName := fmt.Sprintf("%s/%s/%s", ctx, user_id, id+"."+fileExt)

	writer := bucket.Object(objectName).NewWriter(context.Background())
	writer.ContentType = entity.MimeType(fileExt)
	writer.ACL = []storage.ACLRule{
		{Entity: storage.AllUsers, Role: storage.RoleReader},
	}

	if _, err := io.Copy(writer, fileReader); err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	attrs, err := bucket.Object(objectName).Attrs(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get object attributes: %w", err)
	}

	return attrs.MediaLink, nil
}

func DeleteFile(downloadURL string) error {
	firebase, err := config.InitFirebase()
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase: %w", err)
	}

	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to parse download URL: %w", err)
	}

	if parsedURL.Path == "" && parsedURL.RawPath == "" {
		return fmt.Errorf("invalid download URL, no path found")
	}

	path := parsedURL.Path
	if strings.Contains(path, "/o/") {
		path = strings.SplitN(path, "/o/", 2)[1]
		path, err = url.QueryUnescape(path)
		if err != nil {
			return fmt.Errorf("failed to decode object path: %w", err)
		}
	} else {
		return fmt.Errorf("invalid download URL format, missing '/o/' segment")
	}

	bucket := firebase.StorageClient.Bucket(os.Getenv("FIREBASE_BUCKET_NAME"))
	obj := bucket.Object(path)

	_, err = obj.Attrs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get object attributes: %w", err)
	}

	if err := obj.Delete(context.Background()); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}
