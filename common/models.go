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

type LimitlessApiResponse struct {
	Data struct {
		LifeLogs []LimitlessLifeLog `json:"lifeLogs"`
	} `json:"data"`
	Pagination struct {
		NextCursor string `json:"nextCursor"`
	} `json:"pagination"`
}

type LimitlessLifeLog struct {
	ID         string `json:"id"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Title      string `json:"title"`
	Overview   string `json:"overview"`
	Transcript string `json:"transcript"`
}

/*
type BeeParsed struct {
	StartTime          string
	EndTime            string
	ShortSummary       string
	SummaryLines       []string
	TranscriptionLines []string
}

type OmiParsed struct {
	Timestamp       string
	Title           string
	Overview        string
	TranscriptLines []string
}
*/
