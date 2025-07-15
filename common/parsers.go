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
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseBeeFile(path string) (*PendantExport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	var startTime string
	var endTime string
	var shortSummary string
	var summaryLines []string
	var transcriptionLines []string

	section := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Start: "):
			startTime = strings.TrimSpace(strings.TrimPrefix(line, "Start: "))
		case strings.HasPrefix(line, "End: "):
			endTime = strings.TrimSpace(strings.TrimPrefix(line, "End: "))
		case strings.HasPrefix(line, "Short Summary: "):
			shortSummary = strings.TrimSpace(strings.TrimPrefix(line, "Short Summary: "))
		case line == "Summary:":
			section = "summary"
		case line == "Transcription:":
			section = "transcription"
		default:
			if section == "summary" {
				summaryLines = append(summaryLines, line)
			} else if section == "transcription" {
				transcriptionLines = append(transcriptionLines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %v", err)
	}

	return &PendantExport{
		StartTime:  startTime,
		EndTime:    endTime,
		Title:      shortSummary,
		Overview:   strings.Join(summaryLines, "\n"),
		Transcript: strings.Join(transcriptionLines, "\n"),
		Contents: []ContentBlock{
			{Type: "heading1", Content: shortSummary},
			{Type: "heading2", Content: strings.Join(summaryLines, "\n")},
			{Type: "paragraph", Content: strings.Join(transcriptionLines, "\n")},
		},
	}, nil
}

func ParseOmiFile(path string) (*PendantExport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	var timestamp, title, overview string
	var transcriptLines []string
	inTranscript := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Memory from ") {
			timestamp = strings.TrimSpace(strings.TrimPrefix(line, "Memory from "))
		} else if strings.HasPrefix(line, "Title: ") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "Title: "))
		} else if strings.HasPrefix(line, "Overview: ") {
			overview = strings.TrimSpace(strings.TrimPrefix(line, "Overview: "))
		} else if line == "Transcript:" {
			inTranscript = true
		} else if inTranscript {
			transcriptLines = append(transcriptLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %v", err)
	}

	return &PendantExport{
		StartTime:  timestamp,
		Title:      title,
		Overview:   overview,
		Transcript: strings.Join(transcriptLines, "\n"),
		Contents: []ContentBlock{
			{Type: "heading1", Content: title},
			{Type: "heading2", Content: overview},
			{Type: "paragraph", Content: strings.Join(transcriptLines, "\n")},
		},
	}, nil
}
