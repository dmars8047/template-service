package api

import (
	"fmt"

	"github.com/dmars8047/template-service/pkg/templates"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmailTemplateHandler struct {
	getter  EmailTemplateGetter
	creator EmailTemplateCreator
	deleter EmailTemplateDeleter
}

func NewEmailTemplateHandler(getter EmailTemplateGetter, creator EmailTemplateCreator, deleter EmailTemplateDeleter) *EmailTemplateHandler {
	return &EmailTemplateHandler{getter: getter, creator: creator, deleter: deleter}
}

func (handler *EmailTemplateHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/template-service/email-templates", handler.getTemplates)
	router.GET("/api/template-service/email-templates/:template_id", handler.getTemplate)
	router.POST("/api/template-service/email-templates", handler.createTemplate)
	router.DELETE("/api/template-service/email-templates/:template_id", handler.deleteTemplate)
}

// Gets all email templates
func (handler *EmailTemplateHandler) getTemplates(c *gin.Context) {
	// Get the templates from the store
	templates, err := handler.getter.getAllEmailTemplates()

	// If there was an error, return it to the client
	if err != nil {
		c.JSON(500, gin.H{"error": "An unexpected error occurred when retrieveing email templates"})
		return
	}

	// Return the templates to the client
	c.JSON(200, templates)
}

// Gets an email template by id
func (handler *EmailTemplateHandler) getTemplate(c *gin.Context) {
	// Get the requested template id from the route value :template_id
	var requestedTemplateId = c.Param("template_id")

	// Get the template from the store
	template, err := handler.getter.getEmailTemplate(requestedTemplateId)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when retrieveing an emil template with id: %s", requestedTemplateId)})
		return
	}

	// Return the template to the client
	c.JSON(200, template)
}

// Creates an email template
func (handler *EmailTemplateHandler) createTemplate(c *gin.Context) {
	// Get the request body
	var request templates.EmailTemplateCreateRequest
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

	// Check that the html content is not empty
	if request.HtmlContent == "" {
		c.JSON(400, gin.H{"error": "Html content is required"})
		return
	}

	// Make sure that the html content is under 1 megabtye in size
	if len(request.HtmlContent) >= 1048576 {
		c.JSON(400, gin.H{"error": "Html content must be 1 megabyte or under"})
		return
	}

	// Check that the plain text content is not empty
	if request.PlainTextContent == "" {
		c.JSON(400, gin.H{"error": "Plain text content is required"})
		return
	}

	// Make sure that the plain text content is under 1 megabtye in size
	if len(request.PlainTextContent) >= 1048576 {
		c.JSON(400, gin.H{"error": "Plain text content must be 1 megabyte or under"})
		return
	}

	// Check that the subject is not empty
	if request.Subject == "" {
		c.JSON(400, gin.H{"error": "Subject is required"})
		return
	}

	// Make sure that the subject is under 78 characters
	if len(request.Subject) > 78 {
		c.JSON(400, gin.H{"error": "Subject must be under 79 characters"})
		return
	}

	// if the tokens are empty, set it to an empty array
	if request.Tokens == nil {
		request.Tokens = []templates.Token{}
	}

	template := templates.EmailTemplate{
		Id:               uuid.New().String(),
		Name:             request.Name,
		Tokens:           request.Tokens,
		HtmlContent:      request.HtmlContent,
		PlainTextContent: request.PlainTextContent,
		Subject:          request.Subject,
	}

	// Get the template from the store
	err = handler.creator.createEmailTemplate(&template)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNameConflict {
			c.JSON(409, gin.H{"error": "A template with the same name already exists"})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when creating an email template with name: %s", request.Name)})
		return
	}

	// Return the template to the client
	c.JSON(201, template)
}

// Deletes an email template by id
func (handler *EmailTemplateHandler) deleteTemplate(c *gin.Context) {
	// Get the requested template id from the route value :template_id
	var requestedTemplateId = c.Param("template_id")

	// Get the template from the store
	err := handler.deleter.deleteEmailTemplate(requestedTemplateId)

	// If there was an error, return it to the client
	if err != nil {
		if err == ErrNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": fmt.Sprintf("An unexpected error occurred when deleting an email template with id: %s", requestedTemplateId)})
		return
	}

	// Return the template to the client
	c.JSON(204, nil)
}
