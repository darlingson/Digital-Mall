package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Port       string
	SecretKey  string
	DatabaseDSN string
}

func LoadConfig() Config {
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "50051"
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		databaseDSN = "mongodb://localhost:27017" // Default MongoDB URL
	}

	return Config{
		Port:       port,
		SecretKey:  secretKey,
		DatabaseDSN: databaseDSN,
	}
}

func ConnectDB(dsn string) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return client.Database("auth_service")
}
