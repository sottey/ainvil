package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sottey/ainvil/lib/model"
)

// SaveEntry writes an Entry to ./export/YYYY/MM/DD/{entry.ID}.json
// If overwrite is false and the file exists, it skips saving.
func SaveEntry(entry model.Entry, overwrite bool) error {
	t := entry.Timestamp
	dir := filepath.Join("export", t.Format("2006"), t.Format("01"), t.Format("02"))

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filename := filepath.Join(dir, entry.ID+".json")
	if !overwrite {
		if _, err := os.Stat(filename); err == nil {
			// File already exists, skip
			return nil
		}
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
