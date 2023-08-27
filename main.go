package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Set up the MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString).SetAuth(options.Credential{
		Username: mongoUsername,
		Password: mongoPassword,
	}))

	if err != nil {
		log.Fatal(err)
	}

	templateDatabase := client.Database("core")
	templateCollection := templateDatabase.Collection("templates")

	templateStore := NewMongoTemplateStore(templateCollection)

	templateHandler := &TemplateHandler{store: templateStore}

	if err != nil {
		log.Fatal(err)
	}

	// Initialize the router
	router := gin.Default()

	// Add the route handlers
	router.GET("/api/templates", templateHandler.GetAllTemplates)
	router.GET("/api/templates/:template_id", templateHandler.GetTemplate)
	router.POST("/api/templates", templateHandler.CreateTemplate)
	router.DELETE("/api/templates/:template_id", templateHandler.DeleteTemplate)

	router.Run(":8080")
}
