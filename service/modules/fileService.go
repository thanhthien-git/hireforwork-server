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
	"sync"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App
var once sync.Once

func getFirebaseApp() *firebase.App {
	once.Do(func() {
		credentialsFile := os.Getenv("FIREBASE_CREDENTIALS")
		if credentialsFile == "" {
			log.Fatalf("FIREBASE_CREDENTIALS not set in .env")
		}

		opt := option.WithCredentialsFile(credentialsFile)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalf("Error when create new app")
		}
		firebaseApp = app
	})
	return firebaseApp
}

const maxFileSize = 10 * 1024 * 1024 // 10MB

func UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string, contentType string) (string, error) {
	if header.Size > maxFileSize {
		return "", fmt.Errorf("file too large: %d bytes", header.Size)
	}
	app := getFirebaseApp()

	bucketName := os.Getenv("FIRBASE_BUCKET")
	if bucketName == "" {
		log.Fatalf("Bucket not found")
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("Error when open connect to bucket")
	}
	fileName := createFileName(folder, header)
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
		return "", fmt.Errorf("failed to set ACL: %w", err)
	}
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("Error when retrieve object: ", err.Error())
	}
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, attrs.Name)

	return publicURL, nil
}

func UploadResume(file multipart.File, header *multipart.FileHeader, contentType string) (string, error) {
	return UploadFile(context.Background(), file, header, os.Getenv("FIRBASE_BUCKET_RESUME"), contentType)
}

func UploadImage(file multipart.File, header *multipart.FileHeader, contentType string) (string, error) {
	return UploadFile(context.Background(), file, header, os.Getenv("FIRBASE_BUCKET_PICTURE"), contentType)
}

func createFileName(path string, header *multipart.FileHeader) string {
	ext := strings.ToLower(filepath.Ext(header.Filename))

	switch ext {
	case ".png", ".jpeg", ".jpg", ".docx", ".pdf":
		return path + uuid.New().String() + ext
	default:
		return ""
	}
}

type Config struct {
	FirebaseCredentials string
	BucketName          string
	ResumeBucket        string
	PictureBucket       string
}

func NewConfig() (*Config, error) {
	return &Config{
		FirebaseCredentials: os.Getenv("FIREBASE_CREDENTIALS"),
		BucketName:          os.Getenv("FIRBASE_BUCKET"),
		ResumeBucket:        os.Getenv("FIRBASE_BUCKET_RESUME"),
		PictureBucket:       os.Getenv("FIRBASE_BUCKET_PICTURE"),
	}, nil
}
