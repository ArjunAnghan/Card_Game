package config

// Config holds the configuration settings for the application.
// It includes the MongoDB connection URI and the name of the MongoDB database to use.
type Config struct {
	MongoDBURI      string // The URI for connecting to the MongoDB instance
	MongoDBDatabase string // The name of the MongoDB database to use
}

// LoadConfig loads and returns the configuration settings for the application.
// This function initializes and returns a Config struct with hardcoded values.
// You can update the MongoDB URI and database name to match your specific MongoDB setup.
func LoadConfig() *Config {
	return &Config{
		MongoDBURI:      "mongodb://localhost:27017", // Update this to match your MongoDB setup
		MongoDBDatabase: "mydb",                      // Ensure this matches the database name you're trying to use
	}
}
