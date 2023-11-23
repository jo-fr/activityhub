package util

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// TrimStringLength trims a string to a maximum length, but tries to keep the last sentence-ending punctuation. If no sentence-ending punctuation is found, the last character is used.
func TrimStringLength(input string, maxLen int) string {

	if utf8.RuneCountInString(input) <= maxLen {
		return input
	}

	// Trim the input string to maxLen characters
	trimmed := input[:maxLen]

	// Find the last sentence-ending punctuation (period, exclamation mark, or question mark)
	lastSentenceEnd := strings.LastIndexAny(trimmed, ".!?")

	// If no sentence-ending punctuation is found, use the last character
	if lastSentenceEnd == -1 {
		lastSentenceEnd = maxLen - 1
	}

	// Trim the string to the last sentence-ending punctuation and add "..."
	trimmed = trimmed[:lastSentenceEnd+1] + " [...]"

	return trimmed
}

// RemoveHTMLTags removes all HTML tags from a string.
func RemoveHTMLTags(input string) string {
	// Regular expression pattern to match HTML tags
	htmlTagPattern := "<[^>]*>"

	// Replace HTML tags with an empty string
	regex := regexp.MustCompile(htmlTagPattern)
	cleanedString := regex.ReplaceAllString(input, "")

	return cleanedString
}
