package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DatabaseName string

func ConnectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	DatabaseName = os.Getenv("MONGODB_DATABASE")

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo not reachable:", err)
	}

	Client = client
	log.Println("MongoDB connected")
}

func GetCollection(name string) *mongo.Collection {
	return Client.Database(DatabaseName).Collection(name)
}
