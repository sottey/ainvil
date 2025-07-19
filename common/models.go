/*
Copyright © 2025 sottey

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package common

import (
	"encoding/json"
	"time"
)

type PendantExport struct {
	ID            string          `json:"id"`
	SourceType    string          `json:"sourceType"`
	StartTime     string          `json:"startTime"`
	EndTime       string          `json:"endTime"`
	Title         string          `json:"title"`
	Overview      string          `json:"overview"`
	Transcript    string          `json:"transcript"`
	Contents      []ContentEntry  `json:"contents"`
	Markdown      string          `json:"markdown,omitempty"`
	IsStarred     bool            `json:"isStarred,omitempty"`
	UpdatedAt     string          `json:"updatedAt,omitempty"`
	CreatedAt     string          `json:"createdAt,omitempty"`
	ExportDate    string          `json:"exportDate"`
	ExportVersion string          `json:"exportVersion"`
	SourceFile    string          `json:"sourceFile"`
	DeviceType    string          `json:"deviceType,omitempty"`
	Latitude      string          `json:"latitude,omitempty"`
	Longitude     string          `json:"longitude,omitempty"`
	Address       string          `json:"address,omitempty"`
	Raw           json.RawMessage `json:"raw"`
}

type LimitlessResponse struct {
	Data []LimitlessLifelog `json:"data"`
	Meta Meta               `json:"meta"`
}

type LimitlessLifelog struct {
	ID         string         `json:"id"`
	StartTime  string         `json:"startTime"`
	EndTime    string         `json:"endTime"`
	Markdown   string         `json:"markdown"`
	UpdatedAt  string         `json:"updatedAt"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	IsStarred  bool           `json:"isStarred"`
	Transcript string         `json:"transcript"`
	Contents   []ContentEntry `json:"contents"`
	Text       string         `json:"text"`
	Type       string         `json:"type"`
	Location   Location       `json:"location"`
	Tags       []string       `json:"tags"`
	Source     Source         `json:"source"`
	Raw        json.RawMessage
}

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Address   string `json:"address"`
}

type Source struct {
	DeviceType string `json:"deviceType"`
	SourceType string `json:"sourceType"`
}

type Meta struct {
	Lifelogs MetaLifelogs `json:"lifelogs"`
}

type MetaLifelogs struct {
	NextCursor string `json:"nextCursor"`
}

type BeeParsed struct {
	StartTime          string
	EndTime            string
	DeviceType         string
	ShortSummary       string
	SummaryLines       []string
	TranscriptionLines []string
	Latitude           string
	Longitude          string
	Address            string
}

type AinvilExport struct {
	ID            string         `json:"id"`
	SourceType    string         `json:"sourceType"`
	Title         string         `json:"title"`
	StartTime     time.Time      `json:"startTime"`
	EndTime       time.Time      `json:"endTime"`
	Overview      string         `json:"overview"`
	Transcript    string         `json:"transcript"`
	Contents      []ContentEntry `json:"contents"`
	ExportDate    time.Time      `json:"exportDate"`
	ExportVersion string         `json:"exportVersion"`
	SourceFile    string         `json:"sourceFile"`
	Raw           any            `json:"raw"`
	Markdown      string         `json:"markdown,omitempty"`
}

type ContentEntry struct {
	Type    string `json:"type"`
	Content string `json:"content"`

	SpeakerName       string `json:"speakerName,omitempty"`
	SpeakerIdentifier string `json:"speakerIdentifier,omitempty"`
	StartTime         string `json:"startTime,omitempty"`
	EndTime           string `json:"endTime,omitempty"`
	StartOffsetMs     int    `json:"startOffsetMs,omitempty"`
	EndOffsetMs       int    `json:"endOffsetMs,omitempty"`
}
