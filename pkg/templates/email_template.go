package templates

import "log"

// EmailTemplate is a struct that represents an email template.
type EmailTemplate struct {
	Id               string  `json:"id" bson:"_id"`
	Name             string  `json:"name"`
	HtmlContent      string  `json:"html_content" bson:"htmlContent"`
	PlainTextContent string  `json:"plain_text_content" bson:"plainTextContent"`
	Subject          string  `json:"subject"`
	Tokens           []Token `json:"tokens"`
}

// Method for applying substitutions to the email template.
func (emailTemplate *EmailTemplate) ApplySubstitutions(substitutions map[string]string) error {
	var err error

	emailTemplate.HtmlContent, err = applySubstitutions(emailTemplate.HtmlContent, emailTemplate.Tokens, substitutions)
	if err != nil {
		log.Printf("error applying substitutions to html content to email template, ID: %s, error: %s\n", emailTemplate.Id, err)
		return err
	}

	emailTemplate.PlainTextContent, err = applySubstitutions(emailTemplate.PlainTextContent, emailTemplate.Tokens, substitutions)
	if err != nil {
		log.Printf("error applying substitutions to plain text content to email template, ID: %s, error: %s\n", emailTemplate.Id, err)
		return err
	}

	emailTemplate.Subject, err = applySubstitutions(emailTemplate.Subject, emailTemplate.Tokens, substitutions)
	if err != nil {
		log.Printf("error applying substitutions to subject to email template, ID: %s, error: %s\n", emailTemplate.Id, err)
		return err
	}

	return nil
}
