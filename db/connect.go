package dbHelper

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, context.Context, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	clientOption := options.Client().ApplyURI(os.Getenv("DATABASE_CONNECTION"))

	//time > 10s cancel connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//connect to the db
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, nil, err
	}
	//notify when success
	log.Println("Connect to MongoDB")
	return client, ctx, nil
}

func ValidateError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCollection(ctx context.Context, collectionName string, client *mongo.Client) *mongo.Collection {
	Collection := client.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	return Collection
}
