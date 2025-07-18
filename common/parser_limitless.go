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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func toJSON(v any) []byte {
	b, _ := json.MarshalIndent(v, "", "  ")
	return b
}

func ParseLimitlessData(apiKey, apiURL, start, outputDir string) error {
	if apiKey == "" || apiURL == "" {
		return errors.New("missing --token or --url")
	}

	if start == "" {
		fmt.Println("No --start provided. Trying to detect most recent saved date in output...")
		var mostRecent time.Time
		filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasPrefix(filepath.Base(path), "limitless_") || !strings.HasSuffix(path, ".json") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			var exp PendantExport
			if err := json.Unmarshal(data, &exp); err == nil && exp.UpdatedAt != "" {
				if t, err := time.Parse(time.RFC3339, exp.UpdatedAt); err == nil && t.After(mostRecent) {
					mostRecent = t
				}
			}
			return nil
		})
		if !mostRecent.IsZero() {
			start = mostRecent.Format("2006-01-02")
			fmt.Println("Using latest discovered date:", start)
		} else {
			fmt.Println("No previous files found. Proceeding with no start date filter.")
		}
	}

	page := 1
	cursor := ""

	for {
		reqURL := fmt.Sprintf("%s?limit=100", apiURL)
		if cursor != "" {
			reqURL += "&cursor=" + cursor
		} else if start != "" {
			reqURL += "&start=" + start
		}
		fmt.Printf("Fetching page %d...\n", page)
		fmt.Println("Request URL:", reqURL)

		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header.Set("X-API-Key", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("request error: %w", err)
		}
		if resp.StatusCode == 429 {
			fmt.Println("Rate limit hit. Sleeping 20s then retrying...")
			time.Sleep(20 * time.Second)
			resp, err = http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode == 429 {
				return fmt.Errorf("still failing after retry: %v", resp.Status)
			}
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 401 || strings.Contains(string(body), "Unauthorized") {
			return errors.New("unauthorized: check API token")
		}

		var rawResp map[string]interface{}
		if err := json.Unmarshal(body, &rawResp); err != nil {
			fmt.Println("JSON error:", err)
			return err
		}

		dataList := []json.RawMessage{}
		if dataNode, ok := rawResp["data"].(map[string]interface{}); ok {
			if lifelogs, ok := dataNode["lifelogs"].([]interface{}); ok {
				for _, item := range lifelogs {
					if m, err := json.Marshal(item); err == nil {
						dataList = append(dataList, m)
					}
				}
			} else {
				fmt.Println("Warning: 'data.lifeLogs' is missing or not an array")
			}
		} else {
			fmt.Println("Warning: 'data' node missing or not a map")
		}

		if len(dataList) == 0 {
			fmt.Println("No lifelogs found.")
			break
		}
		fmt.Printf("Found %d lifelogs\n", len(dataList))

		for _, raw := range dataList {
			var item LimitlessLifeLog
			if err := json.Unmarshal(raw, &item); err != nil {
				fmt.Println("Skipping malformed lifelog:", err)
				continue
			}
			item.Raw = raw

			export := PendantExport{
				ID:            item.ID,
				SourceType:    "limitless",
				StartTime:     item.StartTime,
				EndTime:       item.UpdatedAt,
				Title:         item.Title,
				Overview:      item.Summary,
				Transcript:    item.Markdown,
				Contents:      item.Contents,
				UpdatedAt:     item.UpdatedAt,
				IsStarred:     item.IsStarred,
				ExportDate:    time.Now().UTC().Format(time.RFC3339),
				ExportVersion: "Ainvil 2.0.0",
				SourceFile:    "limitlessAPI",
				Raw:           raw,
			}

			// Determine output directory
			t, err := time.Parse(time.RFC3339, export.StartTime)
			if err != nil {
				fmt.Printf("Warning: bad StartTime (%s) on lifelog %s, using current time\n", export.StartTime, export.ID)
				t = time.Now().UTC()
			}
			outDir := filepath.Join(outputDir,
				fmt.Sprintf("%04d", t.Year()),
				fmt.Sprintf("%02d", t.Month()),
				fmt.Sprintf("%02d", t.Day()))
			os.MkdirAll(outDir, 0755)

			outPath := filepath.Join(outDir, fmt.Sprintf("limitless_%s.json", export.ID))
			if err := os.WriteFile(outPath, toJSON(export), 0644); err != nil {
				fmt.Println("Failed to save", export.ID, ":", err)
			} else {
				fmt.Println("Saved", export.ID)
			}
		}

		// Detect nextCursor from meta.lifelog.nextCursor
		meta, ok := rawResp["meta"].(map[string]interface{})
		if !ok {
			fmt.Println("No meta node. Ending pagination.")
			break
		}
		lifelogsMeta, ok := meta["lifelogs"].(map[string]interface{})
		if !ok {
			fmt.Println("No lifelog node in meta. Ending pagination.")
			break
		}
		nextCursorRaw, ok := lifelogsMeta["nextCursor"].(string)
		if !ok || nextCursorRaw == "" {
			fmt.Println("No nextCursor. Ending pagination.")
			break
		}

		cursor = nextCursorRaw
		page++
	}

	return nil
}
