package helpers

import (
	"fmt"
	"os"
)

type Config struct {
	DBConnectionString string
	DBName             string
	CollectionName     string
}

func ConfigSetup() (Config, error) {
	dsn := fmt.Sprintf("mongodb://%s:27017", os.Getenv("MONGO_HOST"))
	config := Config{
		DBConnectionString: dsn,
		DBName:             os.Getenv("MONGO_NAME"),
		CollectionName:     os.Getenv("MONGO_COLLECTION"),
	}

	if config.DBConnectionString == "" {
		return config, fmt.Errorf("DB_CONNECTION_STRING not set")
	}
	if config.DBName == "" {
		return config, fmt.Errorf("DB_NAME not set")
	}
	if config.CollectionName == "" {
		return config, fmt.Errorf("COLLECTION_NAME not set")
	}

	fmt.Println(config.DBConnectionString)
	return config, nil
}
