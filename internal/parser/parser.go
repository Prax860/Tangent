package parser

import (
	"regexp"
	"strings"
)

// Parse normalizes user input.
//
// Example:
//
//	"  Create   Virtual Environment!!!  "
//
// becomes
//
//	"create virtual environment"
func Parse(input string) string {

	// Convert to lowercase
	input = strings.ToLower(input)

	// Remove leading and trailing spaces
	input = strings.TrimSpace(input)

	// Remove punctuation except
	// letters, digits, underscore, slash, dot and hyphen
	re := regexp.MustCompile(`[^\w\s./-]`)
	input = re.ReplaceAllString(input, "")

	// Collapse multiple spaces into one
	re = regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	return input
}