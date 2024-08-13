package config

type Config struct {
	MongoDBURI      string
	MongoDBDatabase string
}

func LoadConfig() *Config {
	return &Config{
		MongoDBURI:      "mongodb://localhost:27017", // Update this to match your MongoDB setup
		MongoDBDatabase: "mydb",                      // Ensure this matches the database name you're trying to use
	}
}
