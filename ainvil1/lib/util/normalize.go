package util

import (
	"strings"
	"time"
)

// NormalizeTimestamp parses and safely returns a time.Time from a string.
// Falls back to zero time if invalid.
func NormalizeTimestamp(s string) time.Time {
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

// NormalizeText cleans up strings from AI output or transcripts.
func NormalizeText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\u00A0", " ") // replace non-breaking spaces
	return s
}

// TruncateText returns the first n characters (with ellipsis if needed).
func TruncateText(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "…"
}
