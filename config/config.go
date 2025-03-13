package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseName       string
	SecretKey          string
	MongoUrl           string
	FirebaseBucket     string
	FirebaseCredential string
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system env vars")
		}

		dbName := os.Getenv("DATABASE_NAME")
		secretKey := os.Getenv("SECRET_KEY")
		mongoUrl := os.Getenv("DATABASE_CONNECTION")
		firebaseBucket := os.Getenv("FIREBASE_BUCKET")
		firebaseCredential := os.Getenv("FIREBASE_CREDENTIALS")

		instance = &Config{
			DatabaseName:       dbName,
			MongoUrl:           mongoUrl,
			SecretKey:          secretKey,
			FirebaseBucket:     firebaseBucket,
			FirebaseCredential: firebaseCredential,
		}
	})
	return instance
}
