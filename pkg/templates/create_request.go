package templates

// EmailTemplateCreateRequest is a struct that represents a request to create an email template.
type EmailTemplateCreateRequest struct {
	Name             string  `json:"name"`
	Tokens           []Token `json:"tokens"`
	HtmlContent      string  `json:"html_content"`
	PlainTextContent string  `json:"plain_text_content"`
	Subject          string  `json:"subject"`
}
