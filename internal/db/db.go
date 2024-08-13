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

// ConnectDB establishes a connection to the MongoDB instance using the provided configuration settings.
// It initializes the global MongoDB client and the game database instance.
func ConnectDB(cfg *config.Config) {
	// Configure MongoDB client options with the provided URI
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)

	var err error
	// Create a new MongoDB client
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		// Log and exit if the client creation fails
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Set a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Attempting to connect to MongoDB...")
	// Attempt to connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		// Log and exit if the connection fails
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Pinging MongoDB...")
	// Ping MongoDB to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		// Log and exit if the ping fails
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connected successfully!")

	// Initialize the game database
	gameDB = client.Database(cfg.MongoDBDatabase)
	if gameDB == nil {
		// Log and exit if the database initialization fails
		log.Fatal("Database initialization failed. gameDB is nil.")
	} else {
		log.Println("Database initialized successfully!")
	}
}

// GetCollection returns a reference to a MongoDB collection in the game database.
// It ensures that the database connection is established before accessing collections.
func GetCollection(collectionName string) *mongo.Collection {
	if gameDB == nil {
		// Log and exit if the database connection is nil
		log.Fatal("Database connection is nil. Ensure ConnectDB is called before accessing collections.")
	}
	// Return the requested collection
	return gameDB.Collection(collectionName)
}

// DisconnectDB disconnects from the MongoDB instance and cleans up the client resources.
// It checks if the client is not nil before attempting to disconnect.
func DisconnectDB() {
	if client == nil {
		log.Println("MongoDB client is nil. Skipping disconnect.")
		return
	}

	// Attempt to disconnect from MongoDB
	if err := client.Disconnect(context.Background()); err != nil {
		// Log and exit if the disconnection fails
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB!")
}
