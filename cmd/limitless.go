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
	"strings"
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

		// Auto-detect start date if none specified
		if startStr == "" {
			fmt.Println("No --start provided. Scanning existing output to determine latest saved date...")
			detected, err := findMostRecentSavedDate(outRoot)
			if err != nil {
				fmt.Printf("Error detecting last saved date: %v\n", err)
				os.Exit(1)
			}
			startStr = detected
			fmt.Printf("Auto-selected --start: %s\n", startStr)
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
			//fmt.Printf("URL: '%v', Key: '%v', cursor: '%v'", apiURL, apiKey, cursor)
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

			// sort logs by StartTime
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].StartTime < logs[j].StartTime
			})

			savedThisPage := 0
			for _, ll := range logs {
				startT, err := parseTimeISO(ll.StartTime)
				if err != nil {
					fmt.Printf("Skipping %q: invalid timestamp\n", ll.StartTime)
					continue
				}

				// Apply date filtering
				if !startT.IsZero() && startT.Before(startDate) {
					continue
				}
				if !endDate.IsZero() && startT.After(endDate) {
					continue
				}

				// Build PendantExport
				export := &common.PendantExport{
					ID:            ll.ID,
					SourceType:    "limitless",
					StartTime:     ll.StartTime,
					EndTime:       ll.EndTime,
					Title:         ll.Title,
					Overview:      ll.Overview,
					Transcript:    ll.Transcript,
					ExportDate:    time.Now().UTC().Format(time.RFC3339),
					ExportVersion: ainvilVersion,
					SourceFile:    "limitlessAPI",
					Contents: []common.ContentBlock{
						{Type: "heading1", Content: ll.Title},
						{Type: "heading2", Content: ll.Overview},
						{Type: "paragraph", Content: ll.Transcript},
					},
				}

				if err := saveLimitlessExport(outRoot, export); err != nil {
					fmt.Printf("Error saving %s: %v\n", export.ID, err)
				} else {
					fmt.Printf("Saved %s\n", export.ID)
					savedThisPage++
					totalSaved++
				}
			}

			if apiResp.Pagination.NextCursor == "" {
				break
			}
			cursor = apiResp.Pagination.NextCursor
			pageNum++
		}

		fmt.Printf("Done. %d lifelogs saved.\n", totalSaved)
	},
}

func init() {
	rootCmd.AddCommand(limitlessCmd)
	limitlessCmd.Flags().String("token", "", "API token (required)")
	limitlessCmd.Flags().String("url", "", "Base API URL (required)")
	limitlessCmd.Flags().String("start", "", "Start date (RFC3339) or auto-detected if omitted")
	limitlessCmd.Flags().String("end", "", "End date (RFC3339, optional)")
	limitlessCmd.Flags().String("out", "./out", "Output root directory")
}

// --- Auto-detect logic ---
func findMostRecentSavedDate(outRoot string) (string, error) {
	var latest time.Time
	foundAny := false

	err := filepath.Walk(outRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(filepath.Base(path), "limitless_") || !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		var export common.PendantExport
		if json.NewDecoder(f).Decode(&export) != nil {
			return nil
		}

		t, err := time.Parse(time.RFC3339, export.StartTime)
		if err != nil {
			return nil
		}

		if !foundAny || t.After(latest) {
			latest = t
			foundAny = true
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if !foundAny {
		return "", fmt.Errorf("no existing limitless_*.json files found in %s", outRoot)
	}

	return latest.Format(time.RFC3339), nil
}

// --- existing helper code (unchanged) ---

func parseDateFlag(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

func parseTimeISO(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func saveLimitlessExport(outRoot string, export *common.PendantExport) error {
	t, err := time.Parse(time.RFC3339, export.StartTime)
	if err != nil {
		t = time.Now().UTC()
	}

	outDir := filepath.Join(
		outRoot,
		"limitless",
		fmt.Sprintf("%04d", t.Year()),
		fmt.Sprintf("%02d", t.Month()),
		fmt.Sprintf("%02d", t.Day()),
	)

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	outPath := filepath.Join(outDir, fmt.Sprintf("limitless_%s.json", export.ID))

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}

// --- API fetch (unchanged) ---

type ApiResponse struct {
	Data struct {
		LifeLogs []LifeLog `json:"lifeLogs"`
	} `json:"data"`
	Pagination struct {
		NextCursor string `json:"nextCursor"`
	} `json:"pagination"`
}

type LifeLog struct {
	ID         string `json:"id"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Title      string `json:"title"`
	Overview   string `json:"overview"`
	Transcript string `json:"transcript"`
}

func fetchPage(apiURL, apiKey, cursor string) (*ApiResponse, int, error) {
	client := &http.Client{}
	reqURL, _ := url.Parse(apiURL)
	q := reqURL.Query()
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", reqURL.String(), nil)
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}

	var result ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decoding response: %w", err)
	}

	return &result, resp.StatusCode, nil
}
