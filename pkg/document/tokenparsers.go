package document

import (
	port "github.com/reiver/go-porterstemmer"
	"strings"
)

var punctuation = []string{"!", ".", "?", `"`, `'`, ","}

//We don't want all the punctuation to be in use, remove the one we don't need.
func RemovePunctuation(token string) string {
	for _, value := range punctuation {
		if value == token {
			return ""
		}
	}
	return token
}

func PorterStemmer(token string) string {
	return port.StemString(token)
}

func ToLower(token string) string {
	return strings.ToLower(token)
}
