package complianceai

import "regexp"

var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\b\d{17}[0-9Xx]\b`),
	regexp.MustCompile(`\b1[3-9]\d{9}\b`),
	regexp.MustCompile(`\b\d{16,19}\b`),
}

func RedactSensitive(value string) string {
	for _, pattern := range sensitivePatterns {
		value = pattern.ReplaceAllString(value, "[已脱敏]")
	}
	return value
}
