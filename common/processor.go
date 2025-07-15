package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const ExportVersion = "ainvil 1.0.0"

// ParserFunc defines the signature for a parser that turns a .txt file into PendantExport
type ParserFunc func(path string) (*PendantExport, error)

// ProcessTextExports processes all .txt files in sourceDir using the parser
// and writes them as JSON to outDir/sourceType/yyyy/MM/dd.
func ProcessTextExports(sourceDir, outDir, sourceType string, parser ParserFunc) error {
	if sourceDir == "" {
		return fmt.Errorf("--source is required")
	}

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("error reading source directory: %w", err)
	}

	totalSaved := 0
	for _, entry := range files {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		inputPath := filepath.Join(sourceDir, entry.Name())
		absInputPath, _ := filepath.Abs(inputPath)

		export, err := parser(inputPath)
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		// Fill standard fields
		export.ID = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		export.SourceType = sourceType
		export.ExportDate = time.Now().UTC().Format(time.RFC3339)
		export.ExportVersion = ExportVersion
		export.SourceFile = absInputPath

		// Write it
		if err := saveExport(outDir, export); err != nil {
			fmt.Printf("Error saving %s: %v\n", entry.Name(), err)
		} else {
			fmt.Printf("Saved %s\n", entry.Name())
			totalSaved++
		}
	}

	fmt.Printf("Done. %d memories saved.\n", totalSaved)
	return nil
}

// saveExport writes a PendantExport as JSON to the correct outDir structure.
func saveExport(outRoot string, export *PendantExport) error {
	t, err := time.Parse(time.RFC3339, export.StartTime)
	if err != nil {
		t = time.Now().UTC() // fallback
	}

	outDir := filepath.Join(
		outRoot,
		fmt.Sprintf("%04d", t.Year()),
		fmt.Sprintf("%02d", t.Month()),
		fmt.Sprintf("%02d", t.Day()),
	)

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	outPath := filepath.Join(outDir, export.ID+".json")

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}
