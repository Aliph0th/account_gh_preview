package utils

import "strings"

// FormatTokensWithSeparator splits a number string into groups of 3 digits from right to left,
// separated by spaces. For example: "1234567" becomes "1 234 567"
func FormatTokensWithSeparator(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	var result strings.Builder
	for i := range runes {
		if i > 0 && i%3 == 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(runes[i])
	}

	resultRunes := []rune(result.String())
	for i, j := 0, len(resultRunes)-1; i < j; i, j = i+1, j-1 {
		resultRunes[i], resultRunes[j] = resultRunes[j], resultRunes[i]
	}

	return string(resultRunes)
}
