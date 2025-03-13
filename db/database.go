package db

import (
	"context"
	"hireforwork-server/config"
	"log"
	"sync"

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
		clientOptions := options.Client().ApplyURI(mongoUrl)

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		instance = &DB{Client: client}
		log.Println("MongoDB connected successfully")
	})
	return instance
}

func (d *DB) Close() {
	if d.Client != nil {
		if err := d.Client.Disconnect(context.Background()); err != nil {
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
