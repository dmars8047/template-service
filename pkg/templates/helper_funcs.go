package templates

import (
	"log"
	"sort"
	"strings"
)

func applySubstitutions(content string, tokens []Token, substitutions map[string]string) (string, error) {

	// Iterate through the tokens and replace the toen value with the substitution value.
	// Tokens are replaced in the order of their sequence number.
	// If the token value is not found in the substitutions map, return an error.
	sortBySequenceNumber(tokens)

	for _, token := range tokens {
		if subValue, ok := substitutions[token.Value]; ok {
			content = strings.ReplaceAll(content, token.Value, subValue)
		} else {
			log.Printf("Token not found in substitutions: %s\n", token.Value)
			return "", ErrSubstitutionTokenMismatch
		}
	}

	return content, nil
}

func sortBySequenceNumber(slice []Token) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].SequenceNumber < slice[j].SequenceNumber
	})
}
