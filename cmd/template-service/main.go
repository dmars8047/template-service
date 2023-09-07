package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmars8047/template-service/internal/api"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Starting Template Service...")
	fmt.Printf("Arg count: %d\n", len(os.Args))

	// setup a context with a 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get mongo database connection string from environment variable
	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")

	// If the connection string is empty panic
	if mongoConnectionString == "" {
		panic("MONGO_CONNECTION_STRING environment variable not set")
	}

	// Get the mongo db username from the environment variable
	mongoUsername := os.Getenv("MONGO_USERNAME")

	// If the username is empty panic
	if mongoUsername == "" {
		panic("MONGO_USERNAME environment variable not set")
	}

	// Get the mongo db password from the environment variable
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	// If the password is empty panic
	if mongoPassword == "" {
		panic("MONGO_PASSWORD environment variable not set")
	}

	mongoDatabase := os.Getenv("MONGO_DATABASE")

	if mongoDatabase == "" {
		panic("MONGO_DATABASE environment variable not set")
	}

	mongoCollection := os.Getenv("MONGO_COLLECTION")

	if mongoCollection == "" {
		panic("MONGO_COLLECTION environment variable not set")
	}

	// Set up the MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString).SetAuth(options.Credential{
		Username:   mongoUsername,
		Password:   mongoPassword,
		AuthSource: mongoDatabase,
	}))

	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		panic(err)
	}

	// Ping the MongoDB server to check if the connection was successful
	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB!")

	templateDatabase := client.Database(mongoDatabase)
	templateCollection := templateDatabase.Collection(mongoCollection)

	templateStore := api.NewMongoEmailTemplateStore(templateCollection)

	emailTemplateHandler := api.NewEmailTemplateHandler(templateStore)

	if err != nil {
		log.Fatal(err)
	}

	// Initialize the router
	router := gin.Default()

	// Register the routes
	emailTemplateHandler.RegisterRoutes(router.Group("/api"))

	router.Run(":8080")
}
