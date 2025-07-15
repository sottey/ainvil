package model

import "time"

// Entry represents a normalized lifelog entry from any supported API.
type Entry struct {
	ID        string            `json:"id"`
	Source    string            `json:"source"`             // e.g., "omi", "limitless"
	Timestamp time.Time         `json:"timestamp"`          // normalized start time
	Content   string            `json:"content"`            // raw or formatted text
	Tags      []string          `json:"tags,omitempty"`     // optional tags
	Metadata  map[string]string `json:"metadata,omitempty"` // extra fields, optional
}
