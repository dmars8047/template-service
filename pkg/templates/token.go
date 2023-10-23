package templates

// Token is a struct that represents variable content in a template.
type Token struct {
	// Value is the string value token that will be replaced in the template.
	Value string `json:"value"`
	// SequenceNumber is the order in which the token substitutions will be applied to the template.
	SequenceNumber uint8 `json:"sequence_number" bson:"sequenceNumber"`
}
