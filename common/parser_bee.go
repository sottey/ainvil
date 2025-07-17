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
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func ParseBeeFile(path string) (*PendantExport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	var startTime string
	var endTime string
	var deviceType string
	var shortSummary string
	var summaryLines []string
	var transcriptionLines []string
	var latitude string
	var longitude string
	var address string

	section := ""

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Start Time:"):
			startTime = parseHeaderValue(scanner, line, "Start Time:")
			startTime = parseBeeTimestamp(startTime)

		case strings.HasPrefix(line, "End Time:"):
			endTime = parseHeaderValue(scanner, line, "End Time:")
			endTime = parseBeeTimestamp(endTime)

		case strings.HasPrefix(line, "Device Type:"):
			deviceType = parseHeaderValue(scanner, line, "Device Type:")

		case strings.HasPrefix(line, "Short Summary:"):
			shortSummary = parseHeaderValue(scanner, line, "Short Summary:")

		case line == "Summary:":
			section = "summary"

		case line == "Transcription:":
			section = "transcription"

		case strings.HasPrefix(line, "Primary Location:"):
			section = "location"

		case strings.HasPrefix(line, "Latitude:"):
			latitude = parseHeaderValue(scanner, line, "Latitude:")

		case strings.HasPrefix(line, "Longitude:"):
			longitude = parseHeaderValue(scanner, line, "Longitude:")

		case strings.HasPrefix(line, "bAddress:"):
			address = parseHeaderValue(scanner, line, "bAddress:")

		default:
			switch section {
			case "summary":
				summaryLines = append(summaryLines, line)
			case "transcription":
				transcriptionLines = append(transcriptionLines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %v", err)
	}

	// Construct raw parsed struct
	parsed := struct {
		StartTime          string
		EndTime            string
		DeviceType         string
		ShortSummary       string
		SummaryLines       []string
		TranscriptionLines []string
		Latitude           string
		Longitude          string
		Address            string
	}{
		StartTime:          startTime,
		EndTime:            endTime,
		DeviceType:         deviceType,
		ShortSummary:       shortSummary,
		SummaryLines:       summaryLines,
		TranscriptionLines: transcriptionLines,
		Latitude:           latitude,
		Longitude:          longitude,
		Address:            address,
	}

	rawBytes, _ := json.Marshal(parsed)

	return &PendantExport{
		StartTime:  startTime,
		EndTime:    endTime,
		DeviceType: deviceType,
		Latitude:   latitude,
		Longitude:  longitude,
		Address:    address,
		Title:      shortSummary,
		Overview:   strings.Join(summaryLines, "\n"),
		Transcript: strings.Join(transcriptionLines, "\n"),
		Contents: []ContentEntry{
			{Type: "heading1", Content: shortSummary},
			{Type: "heading2", Content: strings.Join(summaryLines, "\n")},
			{Type: "paragraph", Content: strings.Join(transcriptionLines, "\n")},
		},
		Raw: rawBytes,
	}, nil
}

func parseBeeTimestamp(raw string) string {
	layout := "Jan 2, 2006 at 3:04 PM"
	t, err := time.Parse(layout, raw)
	if err != nil {
		fmt.Printf("Warning: couldn't parse Bee time %q: %v\n", raw, err)
		return raw // preserve original unparsed
	}
	return t.Format(time.RFC3339)
}

func parseHeaderValue(scanner *bufio.Scanner, line, prefix string) string {
	remainder := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if remainder != "" {
		return remainder
	}
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
