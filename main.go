package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	vault "github.com/hashicorp/vault/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getSecret(cfg *Config) (*vault.Secret, error) {
	// 1. สร้าง Vault Client
	config := vault.DefaultConfig()
	config.Address = cfg.VaultAddr

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 2. กำหนด Token
	client.SetToken(cfg.VaultClientID)

	// 3. อ่านข้อมูลจาก Path
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secret, err := client.Logical().Read("kv/mongo")
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	return secret, nil
}

func mongoDBConn(connStr string) {
	clientOptions := options.Client().ApplyURI(connStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func main() {
	const port uint = 8008

	// Load configuration
	cfg := NewConfig()

	// Get secret from Vault
	secret, err := getSecret(cfg)
	if err != nil {
		log.Panic("Error getting secret from Vault:", err)
	}

	if secret == nil || secret.Data == nil {
		log.Panic("No data found at the specified Vault path")
	}

	// Extract MongoDB connection string
	connStr, ok := secret.Data["MONGODB_URI"].(string)
	if !ok || connStr == "" {
		log.Panic("MongoDB connection string not found in Vault secret")
	}

	// Connect to MongoDB
	mongoDBConn(connStr)

	// Start server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Panic("Server failed:", err)
	}

	log.Printf("🚀 Application running on http://127.0.0.1:%d", port)
}
