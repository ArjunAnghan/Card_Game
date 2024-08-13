package db

import (
	"context"
	"log"
	"my-card-game/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	gameDB *mongo.Database
)

func ConnectDB(cfg *config.Config) {
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Attempting to connect to MongoDB...")
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Pinging MongoDB...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connected successfully!")

	gameDB = client.Database(cfg.MongoDBDatabase)
	if gameDB == nil {
		log.Fatal("Database initialization failed. gameDB is nil.")
	} else {
		log.Println("Database initialized successfully!")
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	if gameDB == nil {
		log.Fatal("Database connection is nil. Ensure ConnectDB is called before accessing collections.")
	}
	return gameDB.Collection(collectionName)
}

func DisconnectDB() {
	if client == nil {
		log.Println("MongoDB client is nil. Skipping disconnect.")
		return
	}

	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB!")
}
