package common

import "encoding/json"

type ContentBlock struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type PendantExport struct {
	ID         string `json:"id"`
	SourceType string `json:"sourceType"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime,omitempty"`

	ExportDate    string `json:"exportDate"`
	ExportVersion string `json:"exportVersion"`
	SourceFile    string `json:"sourceFile"`

	Tags      []string `json:"tags,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`

	Title      string         `json:"title,omitempty"`
	Overview   string         `json:"overview,omitempty"`
	Transcript string         `json:"transcript,omitempty"`
	Contents   []ContentBlock `json:"contents,omitempty"`

	Raw json.RawMessage `json:"raw"`
}
