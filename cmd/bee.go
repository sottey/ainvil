package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
)

var beeCmd = &cobra.Command{
	Use:   "bee",
	Short: "Process Bee export text files",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, _ := cmd.Flags().GetString("source")
		outDir, _ := cmd.Flags().GetString("out")

		err := common.ProcessTextExports(sourceDir, outDir, "bee", parseBeeFile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

/*
	var beeCmd = &cobra.Command{
		Use:   "bee",
		Short: "Process Bee export text files",
		Run: func(cmd *cobra.Command, args []string) {
			sourceDir, _ := cmd.Flags().GetString("source")
			outDir, _ := cmd.Flags().GetString("out")

			if sourceDir == "" {
				fmt.Println("Error: --source is required")
				cmd.Usage()
				os.Exit(1)
			}

			files, err := os.ReadDir(sourceDir)
			if err != nil {
				fmt.Printf("Error reading source directory: %v\n", err)
				os.Exit(1)
			}

			totalSaved := 0
			for _, entry := range files {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
					continue
				}

				inputPath := filepath.Join(sourceDir, entry.Name())
				absInputPath, err := filepath.Abs(inputPath)
				if err != nil {
					absInputPath = inputPath // fallback
				}

				parsed, err := parseBeeFile(inputPath)
				if err != nil {
					fmt.Printf("Skipping %s: %v\n", entry.Name(), err)
					continue
				}

				rawBytes, _ := json.Marshal(parsed)

				export := common.PendantExport{
					ID:            strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
					SourceType:    "bee",
					StartTime:     parsed.StartTime,
					EndTime:       parsed.EndTime,
					ExportDate:    time.Now().UTC().Format(time.RFC3339),
					ExportVersion: ainvilVersion,
					SourceFile:    absInputPath,
					Title:         parsed.ShortSummary,
					Overview:      strings.Join(parsed.SummaryLines, "\n"),
					Transcript:    strings.Join(parsed.TranscriptionLines, "\n"),
					Contents: []common.ContentBlock{
						{Type: "heading1", Content: parsed.ShortSummary},
						{Type: "heading2", Content: strings.Join(parsed.SummaryLines, "\n")},
						{Type: "paragraph", Content: strings.Join(parsed.TranscriptionLines, "\n")},
					},
					Raw: rawBytes,
				}

				if err := saveBeeExport(outDir, export, entry.Name(), "bee"); err != nil {
					fmt.Printf("Error saving %s: %v\n", entry.Name(), err)
				} else {
					fmt.Printf("Saved %s\n", entry.Name())
					totalSaved++
				}
			}

			fmt.Printf("Done. %d memories saved.\n", totalSaved)
		},
	}
*/
func init() {
	rootCmd.AddCommand(beeCmd)
	beeCmd.Flags().String("source", "", "Directory containing Bee .txt files (required)")
	beeCmd.Flags().String("out", "./out", "Output root directory")
}

type beeParsed struct {
	StartTime          string
	EndTime            string
	ShortSummary       string
	SummaryLines       []string
	TranscriptionLines []string
}

func parseBeeFile(path string) (*common.PendantExport, error) {
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

	return &common.PendantExport{
		StartTime:  startTime,
		EndTime:    endTime,
		Title:      shortSummary,
		Overview:   strings.Join(summaryLines, "\n"),
		Transcript: strings.Join(transcriptionLines, "\n"),
		Contents: []common.ContentBlock{
			{Type: "heading1", Content: shortSummary},
			{Type: "heading2", Content: strings.Join(summaryLines, "\n")},
			{Type: "paragraph", Content: strings.Join(transcriptionLines, "\n")},
		},
	}, nil
}

func saveBeeExport(outRoot string, export common.PendantExport, filename string, prefix string) error {
	// Parse StartTime for folder layout
	t, err := time.Parse(time.RFC3339, export.StartTime)
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

	outName := fmt.Sprintf("%s_%s.json", prefix, strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)))
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
