package db

import (
	"context"
	"hireforwork-server/config"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
}

var instance *DB
var once sync.Once

func GetInstance() *DB {
	once.Do(func() {
		mongoUrl := config.GetInstance().MongoUrl
		clientOptions := options.Client().ApplyURI(mongoUrl).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}

		instance = &DB{Client: client}
		log.Println("MongoDB connected successfully")
	})
	return instance
}

func (d *DB) Close() {
	if d.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := d.Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}
}

func (d *DB) GetCollection(name string) *mongo.Collection {
	db := GetInstance()
	return db.Database(config.GetInstance().DatabaseName).Collection(name)
}

func (d *DB) GetCollections(names []string) []*mongo.Collection {
	collections := make([]*mongo.Collection, 0, len(names))

	for _, name := range names {
		collections = append(collections, d.GetCollection(name))
	}

	return collections
}
