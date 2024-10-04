package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"serverdb/ent"
	"serverdb/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func implementServer(client *ent.Client) {
	os.Setenv("GIN_MODE", "release")
	gin.SetMode(gin.ReleaseMode)

	errEnv := godotenv.Load()
	if errEnv != nil {
		fmt.Printf("Can't load ENV file: %s\n", errEnv)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
		fmt.Printf("PORT not found in ENV, using default PORT: %s\n", PORT)
	} else {
		fmt.Printf("PORT: %s\n", PORT)
	}

	router := gin.Default()

	handler.RegisterRoute(router, client)

	err := router.Run(fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	router.Run(fmt.Sprintf("localhost:%s", PORT))
}

func implementDatabase() (*ent.Client, error) {
	errEnv := godotenv.Load()

	// Checking .env file is exist!
	if errEnv != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get DatabaseURL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// Initialize Ent Client
	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ent Connect Database Failed!")
	}

	// Set *Database and catch an *Error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close() // optimize the database for instance => | Memory Leaks | Connection Limit | Database Timeout (postgres) |

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal("Failed to create schema resources: %w", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database")
	return client, nil
}

func main() {

	dbclient, err := implementDatabase()
	if err != nil {
		log.Fatal("Faild to implement database %w", err)
	}

	defer implementServer(dbclient) // Latest Load (Avoid server load first may cause to not running database)

}
