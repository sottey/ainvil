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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ParserFunc func(path string) (*PendantExport, error)

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
		export.ExportVersion = GetVersion()
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
