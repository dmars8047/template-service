package templates

import "errors"

// ErrTemplateNotFound is an error that is returned when a template is not found.
var ErrTemplateNotFound = errors.New("template not found")

// ErrCouldNotRetrieveTemplate is an error that is returned when a template could not be retrieved.
var ErrCouldNotRetrieveTemplate = errors.New("an error occurred while retrieving the template")

// ErrCouldNotReadTemplate is an error that is returned when a template could not be read.
var ErrCouldNotReadTemplate = errors.New("an error occurred while parsing the template content")

// ErrSubstitutionTokenMismatch is an error that is returned when a substitution token mismatch occurs.
var ErrSubstitutionTokenMismatch = errors.New("substitution token mismatch when applying substitutions, all tokens must be substituted")
