package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
)

var omiCmd = &cobra.Command{
	Use:   "omi",
	Short: "Process Omi export text files",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, _ := cmd.Flags().GetString("source")
		outDir, _ := cmd.Flags().GetString("out")

		err := common.ProcessTextExports(sourceDir, outDir, "omi", parseOmiFile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(omiCmd)
	omiCmd.Flags().String("source", "", "Directory containing Omi .txt files (required)")
	omiCmd.Flags().String("out", "./out", "Output root directory")
}

// Helper struct for parsing original Omi fields
type omiParsed struct {
	Timestamp       string
	Title           string
	Overview        string
	TranscriptLines []string
}

func parseOmiFile(path string) (*common.PendantExport, error) {
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

	// Build *common.PendantExport
	return &common.PendantExport{
		StartTime:  timestamp,
		Title:      title,
		Overview:   overview,
		Transcript: strings.Join(transcriptLines, "\n"),
		Contents: []common.ContentBlock{
			{Type: "heading1", Content: title},
			{Type: "heading2", Content: overview},
			{Type: "paragraph", Content: strings.Join(transcriptLines, "\n")},
		},
	}, nil
}

/*
func saveOmiExport(outRoot string, export common.PendantExport, filename string, prefix string) error {
	// Use startTime for folder layout
	t, err := time.Parse(time.RFC3339, export.StartTime)
	if err != nil {
		// Try RFC3339Nano fallback for Omi
		t, err = time.Parse(time.RFC3339Nano, export.StartTime)
		if err != nil {
			return fmt.Errorf("invalid startTime %q: %v", export.StartTime, err)
		}
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
*/
