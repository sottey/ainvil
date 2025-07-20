package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Lifelog struct {
	ID            string  `json:"id"`
	SourceType    string  `json:"sourceType"`
	DeviceType    string  `json:"deviceType"`
	Title         string  `json:"title"`
	Overview      string  `json:"overview"`
	StartTimeRaw  string  `json:"startTime"`
	EndTimeRaw    string  `json:"endTime"`
	ExportDateRaw string  `json:"exportDate"`
	Latitude      float64 `json:"latitude,string"`
	Longitude     float64 `json:"longitude,string"`
	Address       string  `json:"address"`
	Transcript    string  `json:"transcript"`
	ExportVersion string  `json:"exportVersion"`
	FilePath      string
	RawJSON       string
	StartTime     time.Time
	EndTime       time.Time
	ExportDate    time.Time
}

func parseFlexibleTime(fieldName string, raw string) *time.Time {
	str := strings.TrimSpace(raw)
	if str == "" {
		return nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02", // fallback format
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, str); err == nil {
			return &t
		}
	}
	fmt.Printf("Could not parse %s: %q\n", fieldName, str)
	return nil
}

func lifelogExists(db *sql.DB, id string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM lifelogs WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking for duplicate:", id, err)
		return true
	}
	return exists
}

func importFile(path string, db *sql.DB) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var log Lifelog
	if err := json.Unmarshal(data, &log); err != nil {
		fmt.Println("Skipping invalid JSON:", path, err)
		return nil
	}

	log.RawJSON = string(data)
	log.FilePath = path

	// Fallback: set missing deviceType from raw
	var raw map[string]interface{}
	_ = json.Unmarshal(data, &raw)
	if log.DeviceType == "" {
		if val, ok := raw["deviceType"].(string); ok {
			log.DeviceType = val
		} else if rawInner, ok := raw["raw"].(map[string]interface{}); ok {
			if val, ok := rawInner["DeviceType"].(string); ok {
				log.DeviceType = val
			}
		}
		if log.DeviceType == "" {
			log.DeviceType = log.SourceType
		}
	}

	// Parse start/end/exportDate safely
	start := parseFlexibleTime("startTime", log.StartTimeRaw)
	end := parseFlexibleTime("endTime", log.EndTimeRaw)

	if start == nil && end == nil {
		fmt.Println("Skipping due to missing start/end time:", path)
		return nil
	} else if start == nil {
		start = end
	} else if end == nil {
		end = start
	}
	log.StartTime = *start
	log.EndTime = *end

	if ed := parseFlexibleTime("exportDate", log.ExportDateRaw); ed != nil {
		log.ExportDate = *ed
	}

	// Skip duplicates
	if lifelogExists(db, log.ID) {
		fmt.Println("Skipping duplicate:", log.ID)
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO lifelogs (
			id, source_type, device_type, title, overview,
			start_time, end_time, latitude, longitude, address,
			transcript, export_date, export_version, file_path, raw_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.ID, log.SourceType, log.DeviceType, log.Title, log.Overview,
		log.StartTime, log.EndTime, log.Latitude, log.Longitude, log.Address,
		log.Transcript, log.ExportDate, log.ExportVersion, log.FilePath, log.RawJSON)

	if err != nil {
		fmt.Println("Insert failed:", log.ID, err)
	}
	return nil
}

func ensureSchema(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS lifelogs (
			id TEXT PRIMARY KEY,
			source_type TEXT,
			device_type TEXT,
			title TEXT,
			overview TEXT,
			start_time DATETIME,
			end_time DATETIME,
			latitude REAL,
			longitude REAL,
			address TEXT,
			transcript TEXT,
			export_date DATETIME,
			export_version TEXT,
			file_path TEXT,
			raw_json TEXT
		)`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ainvil_to_db <output-root-dir> <output.sqlite>")
		return
	}

	root := os.Args[1]
	dbfile := os.Args[2]

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ensureSchema(db)

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		return importFile(path, db)
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Import complete.")
}
