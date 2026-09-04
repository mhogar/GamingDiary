package data

import (
	"strings"
)

func Capitalize(base string) string {
	tokens := strings.Split(strings.ToLower(base), " ")
	for i, token := range tokens {
		tokens[i] = capitalizeWord(token, i == 0)
	}
	return strings.Join(tokens, " ")
}

func capitalizeWord(word string, first bool) string {
	if len(word) == 0 {
		return ""
	}

	if len(word) > 1 {
		return strings.ToTitle(word[:1]) + word[1:]
	}

	if first { // len == 1
		return strings.ToTitle(word)
	} else {
		return word
	}
}
