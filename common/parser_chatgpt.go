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
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var metaRE = regexp.MustCompile(`^(Recorder|Timezone|Start|End):\s*(.+)$`)
var lineRE = regexp.MustCompile(`^\[(\d+)\]\s+(Speaker \d+):\s*(.+)$`)

func ParseChatGPTTranscripts(sourceDir, outDir string) error {
	return filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".txt") {
			return nil
		}
		fmt.Println("Parsing:", path)
		return parseFile(path, outDir)
	})
}

func parseFile(filePath, outDir string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var (
		start, end time.Time
		lines      []ContentEntry
		transcript strings.Builder
	)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := metaRE.FindStringSubmatch(line); len(matches) == 3 {
			switch matches[1] {
			case "Start":
				start, _ = time.Parse("2006-01-02 15:04:05.000", matches[2])
			case "End":
				end, _ = time.Parse("2006-01-02 15:04:05.000", matches[2])
			}
		} else if matches := lineRE.FindStringSubmatch(line); len(matches) == 4 {
			offset := matches[1]
			speaker := matches[2]
			text := matches[3]

			lines = append(lines, ContentEntry{
				Type:        "blockquote",
				Content:     text,
				SpeakerName: speaker,
				StartTime:   offset,
			})
			transcript.WriteString(fmt.Sprintf("[%s] %s: %s\n", offset, speaker, text))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	rawData, _ := json.Marshal(map[string]any{
		"sourceFile": filepath.Base(filePath),
		"exportedAt": time.Now().Format(time.RFC3339),
	})

	export := PendantExport{
		ID:            GenerateID(),
		SourceType:    "chatgpt",
		StartTime:     start.Format(time.RFC3339),
		EndTime:       end.Format(time.RFC3339),
		Transcript:    transcript.String(),
		ExportDate:    time.Now().Format("2006-01-02"),
		ExportVersion: "1.0",
		SourceFile:    filepath.Base(filePath),
		Contents:      lines,
		Raw:           rawData,
	}

	folder := filepath.Join(outDir, fmt.Sprintf("%04d", start.Year()), fmt.Sprintf("%02d", start.Month()), fmt.Sprintf("%02d", start.Day()))
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}

	outputFile := filepath.Join(folder, "chatgpt_"+export.ID+".json")
	fo, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer fo.Close()

	encoder := json.NewEncoder(fo)
	encoder.SetIndent("", "  ")
	return encoder.Encode(export)
}
