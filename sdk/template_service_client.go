package sdk

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type TemplateServiceClient interface {
	GetTemplateById(templateId string) (*Template, error)
}

type HttpTemplateServiceClient struct {
	templateServiceBaseUrl string
}

func NewHttpTemplateServiceClient(templateServiceBaseUrl string) *HttpTemplateServiceClient {
	return &HttpTemplateServiceClient{templateServiceBaseUrl: templateServiceBaseUrl}
}

func (client *HttpTemplateServiceClient) GetTemplateById(templateId string) (*Template, error) {
	// Make an http call to get the template by id
	response, err := http.Get(client.templateServiceBaseUrl + "/templates/" + templateId)

	// If there was an error, return it
	if err != nil {
		return nil, err
	}

	// If the status code is not 200, return an error
	if response.StatusCode == 404 {
		return nil, errors.New("template not found")
	}

	// If the status code is not 200, return an error
	if response.StatusCode != 200 {
		return nil, errors.New("error getting template")
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading template response body")
	}

	// Get the template from the response body
	var template Template

	if err := json.Unmarshal(responseBody, &template); err != nil {
		return nil, errors.New("error unmarshalling template response body")
	}

	return &template, nil
}
