package core

import (
	"regexp"
	"strings"
)

func Parse(input string) string {

	input = strings.ToLower(input)

	input = strings.TrimSpace(input)

	re := regexp.MustCompile(`[^\w\s./-]`)
	input = re.ReplaceAllString(input, "")

	re = regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	return input
}
