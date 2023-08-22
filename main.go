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

	// Set up the MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	templateDatabase := client.Database("template_service")
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
