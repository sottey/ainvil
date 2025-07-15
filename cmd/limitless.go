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
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
)

var limitlessCmd = &cobra.Command{
	Use:   "limitless",
	Short: "Fetch and export from Limitless API to standardized format",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("token")
		apiURL, _ := cmd.Flags().GetString("url")
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		outRoot, _ := cmd.Flags().GetString("out")

		if apiKey == "" || apiURL == "" {
			fmt.Println("Error: --token and --url are required")
			cmd.Usage()
			os.Exit(1)
		}

		startDate, err := parseDateFlag(startStr)
		if err != nil {
			fmt.Println("Error parsing --start:", err)
			os.Exit(1)
		}
		endDate, err := parseDateFlag(endStr)
		if err != nil {
			fmt.Println("Error parsing --end:", err)
			os.Exit(1)
		}

		cursor := ""
		totalSaved := 0
		pageNum := 1
		hadFirst429 := false

		for {
			fmt.Printf("Fetching page %d...\n", pageNum)
			apiResp, status, err := fetchPage(apiURL, apiKey, cursor)
			if err != nil && status == 429 {
				if !hadFirst429 {
					fmt.Println("Received 429 Too Many Requests. Waiting 30 seconds before retrying...")
					time.Sleep(30 * time.Second)
					hadFirst429 = true
					continue
				} else {
					fmt.Println("Received 429 Too Many Requests again. Aborting.")
					os.Exit(1)
				}
			}
			if err != nil {
				fmt.Printf("Error fetching page %d: %v\n", pageNum, err)
				os.Exit(1)
			}

			hadFirst429 = false

			logs := apiResp.Data.LifeLogs
			count := len(logs)
			fmt.Printf("Lifelog count for page %d: %d\n", pageNum, count)

			sort.Slice(logs, func(i, j int) bool {
				return logs[i].StartTime < logs[j].StartTime
			})

			for _, ll := range logs {
				startT, err := parseTimeISO(ll.StartTime)
				if err != nil {
					fmt.Printf("Skipping %q: invalid startTime\n", ll.ID)
					continue
				}

				if !inRange(startT, startDate, endDate) {
					continue
				}

				// Build pendant export
				rawBytes, _ := json.Marshal(ll)
				export := common.PendantExport{
					ID:            ll.ID,
					SourceType:    "limitless",
					StartTime:     ll.StartTime,
					EndTime:       ll.EndTime,
					ExportDate:    time.Now().UTC().Format(time.RFC3339),
					ExportVersion: ainvilVersion,
					SourceFile:    apiURL,
					Title:         ll.Title,
					Overview:      ll.Markdown,
					Contents:      parseLimitlessContents(ll.Contents),
					Raw:           rawBytes,
				}

				if err := saveLimitlessExport(outRoot, export, ll.ID, "limitless"); err != nil {
					fmt.Printf("Error saving %s: %v\n", ll.ID, err)
				} else {
					fmt.Printf("Saved %s\n", ll.ID)
					totalSaved++
				}
			}

			nextCursor := apiResp.Meta.Lifelogs.NextCursor
			if nextCursor == nil || *nextCursor == "" {
				fmt.Println("No more pages.")
				break
			}
			cursor = *nextCursor
			pageNum++
		}

		fmt.Printf("Done. %d lifelogs saved.\n", totalSaved)
	},
}

func init() {
	rootCmd.AddCommand(limitlessCmd)
	limitlessCmd.Flags().String("token", "", "API key (required)")
	limitlessCmd.Flags().String("url", "", "API URL (required)")
	limitlessCmd.Flags().String("start", "", "Start date (optional) in MM-DD-YYYY")
	limitlessCmd.Flags().String("end", "", "End date (optional) in MM-DD-YYYY")
	limitlessCmd.Flags().String("out", "./out", "Output root directory")
}

// Limitless-specific structures
type LifeLog struct {
	ID        string         `json:"id"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
	Contents  []ContentBlock `json:"contents"`
	Title     string         `json:"title"`
	Markdown  string         `json:"markdown"`
	IsStarred bool           `json:"isStarred"`
	UpdatedAt string         `json:"updatedAt"`
}

type ContentBlock struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type APIResponse struct {
	Data struct {
		LifeLogs []LifeLog `json:"lifelogs"`
	} `json:"data"`
	Meta struct {
		Lifelogs struct {
			NextCursor *string `json:"nextCursor"`
			Count      *int    `json:"count"`
		} `json:"lifelogs"`
	} `json:"meta"`
}

func parseDateFlag(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse("01-02-2006", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date format %q (use MM-DD-YYYY)", value)
	}
	return &t, nil
}

func parseTimeISO(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func inRange(t time.Time, start, end *time.Time) bool {
	if start != nil && t.Before(*start) {
		return false
	}
	if end != nil && t.After(*end) {
		return false
	}
	return true
}

func fetchPage(apiURL, apiKey, cursor string) (*APIResponse, int, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid API URL: %v", err)
	}

	q := u.Query()
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %v", err)
	}
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("making request: %v", err)
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status == 429 {
		return nil, status, fmt.Errorf("received 429 Too Many Requests")
	}

	if status != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, status, fmt.Errorf("non-200 response: %d\nBody: %s", status, body)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, status, fmt.Errorf("decoding JSON: %v", err)
	}

	return &apiResp, status, nil
}

func parseLimitlessContents(contents []ContentBlock) []common.ContentBlock {
	var result []common.ContentBlock
	for _, c := range contents {
		result = append(result, common.ContentBlock{
			Type:    c.Type,
			Content: c.Content,
		})
	}
	return result
}

func saveLimitlessExport(outRoot string, export common.PendantExport, id string, prefix string) error {
	t, err := parseTimeISO(export.StartTime)
	if err != nil {
		return fmt.Errorf("invalid startTime %q: %v", export.StartTime, err)
	}

	year := t.Format("2006")
	month := t.Format("01")
	day := t.Format("02")
	dir := filepath.Join(outRoot, year, month, day)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %q: %v", dir, err)
	}

	outName := fmt.Sprintf("%s_%s.json", prefix, id)
	outPath := filepath.Join(dir, outName)

	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("Skipping %s: already exists\n", outPath)
		return nil
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating file %q: %v", outPath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}
