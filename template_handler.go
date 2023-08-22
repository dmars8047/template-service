package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TemplateHandler struct {
	store TemplateStore
}

func (handler *TemplateHandler) GetTemplate(c *gin.Context) {
	// Get the requested template id from the route value :template_id
	var requestedTemplateId = c.Param("template_id")

	// Get the template from the store
	template, err := handler.store.GetTemplateById(requestedTemplateId)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when retrieveing template with id: %s", requestedTemplateId)})
		return
	}

	// Return the template to the client
	c.JSON(200, template)
}

func (handler *TemplateHandler) GetAllTemplates(c *gin.Context) {
	// Get the name query parameter
	var name = c.Query("name")

	// Get the templates from the store
	templates, err := handler.store.GetAllTemplates(name)

	// If there was an error, return it to the client
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when retrieveing templates with name: %s", name)})
		return
	}

	// Return the templates to the client
	c.JSON(200, templates)
}

type TemplateCreateRequest struct {
	Name    string   `json:"name"`
	Tokens  []string `json:"tokens"`
	Content string   `json:"content"`
	Format  string   `json:"format"`
}

func (handler *TemplateHandler) CreateTemplate(c *gin.Context) {
	// Get the request body
	var request TemplateCreateRequest
	err := c.BindJSON(&request)

	// If there was an error, return bad request
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Check that the name is not empty
	if request.Name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	// Check that the content is not empty
	if request.Content == "" {
		c.JSON(400, gin.H{"error": "Content is required"})
		return
	}

	// Check that the format is not empty
	if request.Format == "" {
		c.JSON(400, gin.H{"error": "Format is required"})
		return
	}

	// if the tokens are empty, set it to an empty array
	if request.Tokens == nil {
		request.Tokens = []string{}
	}

	template := Template{
		Id:      uuid.New().String(),
		Name:    request.Name,
		Tokens:  request.Tokens,
		Content: request.Content,
		Format:  request.Format,
	}

	// Get the template from the store
	err = handler.store.CreateTemplate(&template)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNameConflict {
			c.JSON(409, gin.H{"error": "A template with the same name already exists"})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when creating template with name: %s", request.Name)})
		return
	}

	// Return the template to the client
	c.JSON(201, template)
}

func (handler *TemplateHandler) DeleteTemplate(c *gin.Context) {
	// Get the requested template id from the route value :template_id
	var requestedTemplateId = c.Param("template_id")

	// Get the template from the store
	err := handler.store.DeleteTemplate(requestedTemplateId)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when deleting template with id: %s", requestedTemplateId)})
		return
	}

	// Return the template to the client
	c.JSON(204, nil)
}
