package document

import (
	port "github.com/reiver/go-porterstemmer"
	"strings"
)

func PorterStemmer(token string) string {
	return port.StemString(token)
}

func ToLower(token string) string {
	return strings.ToLower(token)
}

func StopList(token string) string {
	return port.StemString(token)
}
