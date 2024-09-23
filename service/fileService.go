package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func InitFireBase() *firebase.App {
	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsFile == "" {
		log.Fatalf("FIREBASE_CREDENTIALS not set in .env")
	}

	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error when create new app")
	}
	return app

}

func UploadFile(file multipart.File, header *multipart.FileHeader, folder string, contentType string) (string, error) {
	app := InitFireBase()

	bucketName := os.Getenv("FIRBASE_BUCKET")
	if bucketName == "" {
		log.Fatalf("Bucket not found")
	}

	ctx := context.Background()

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("Error when open connect to bucket")
	}
	fileName := CreateFileName(folder, header)
	object := bucket.Object(fileName)

	wc := object.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to copy file data: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf(err.Error())
	}
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("Error when retrieve object: ", err.Error())
	}
	return attrs.MediaLink, nil
}

func UploadResume(file multipart.File, header *multipart.FileHeader, contentType string) (string, error) {
	return UploadFile(file, header, os.Getenv("FIRBASE_BUCKET_RESUME"), contentType)
}

func UploadImage(file multipart.File, header *multipart.FileHeader, contentType string) (string, error) {
	return UploadFile(file, header, os.Getenv("FIRBASE_BUCKET_PICTURE"), contentType)
}
func CreateFileName(path string, header *multipart.FileHeader) string {
	ext := strings.ToLower(filepath.Ext(header.Filename))

	switch ext {
	case ".png", ".jpeg", ".jpg", ".docx":
		return path + uuid.New().String() + ext
	default:
		return ""
	}
}
