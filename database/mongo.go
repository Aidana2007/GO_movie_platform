package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	MongoDB := os.Getenv("MONGODB_URI")

	if MongoDB == "" {
		log.Fatal("MONGODB_URI not set")
	}

	clientOptions := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client

}

var Client *mongo.Client = ConnectDB()

func GetCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DatabaseName := os.Getenv("MONGODB_DATABASE")
	collection := Client.Database(DatabaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}
	return collection
}
