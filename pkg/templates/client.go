package templates

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// HttpTemplateServiceClient is a struct that implements the TemplateServiceClient interface. Calls the template service over http.
type HttpTemplateServiceClient struct {
	templateServiceBaseUrl string
}

// NewHttpTemplateServiceClient is a function that returns a new HttpTemplateServiceClient.
func NewHttpTemplateServiceClient(templateServiceBaseUrl string) *HttpTemplateServiceClient {
	return &HttpTemplateServiceClient{templateServiceBaseUrl: templateServiceBaseUrl}
}

// GetEmailTemplateById is a method that returns an email template by id.
func (client *HttpTemplateServiceClient) GetEmailTemplateById(templateId string) (*EmailTemplate, error) {
	// Make an http call to get the template by id
	response, err := http.Get(client.templateServiceBaseUrl + getEmailTemplateByIdEndpointUrlSuffix + templateId)

	// If there was an error, return it
	if err != nil {
		return nil, err
	}

	// If the status code is not 200, return an error
	if response.StatusCode == 404 {
		log.Println("email template not found")
		return nil, ErrTemplateNotFound
	}

	// If the status code is not 200, return an error
	if response.StatusCode != 200 {
		log.Println("error getting email template, non 200 status code")
		return nil, ErrCouldNotRetrieveTemplate
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading email template response body")
		return nil, ErrCouldNotReadTemplate
	}

	// Get the email template from the response body
	var template EmailTemplate

	if err := json.Unmarshal(responseBody, &template); err != nil {
		log.Println("error unmarshalling email template response body")
		return nil, ErrCouldNotReadTemplate
	}

	return &template, nil
}
