package validation

import (
	"strings"
	"unicode"
)

func Required(value, field string) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", invalid("%s is required", field)
	}
	return normalized, nil
}

func Compact(value string) string {
	words := strings.Fields(value)
	return strings.Join(words, " ")
}

func Identifier(value, field string) (string, error) {
	normalized, err := Required(value, field)
	if err != nil {
		return "", err
	}
	if len(normalized) > 80 {
		return "", invalid("%s must be at most 80 characters", field)
	}
	for _, char := range normalized {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-' || char == '_' {
			continue
		}
		return "", invalid("%s contains unsupported characters", field)
	}
	return normalized, nil
}
