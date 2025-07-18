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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const currentVersion = "Ainvil 2.1.0"

func FindMostRecentSavedDate(outRoot string) (string, error) {
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

		var export PendantExport
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

func ParseDateFlag(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

func ParseTimeISO(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func GetVersion() string {
	return currentVersion
}

func GenerateID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		// fallback if needed
		return "fallbackid"
	}
	return hex.EncodeToString(b)
}

func AddCommonFileFlags(cmd *cobra.Command) {
	cmd.Flags().String("source", "", "Directory containing input files")
	viper.BindPFlags(cmd.Flags())
}

func AddCommonAPIFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("token", "t", "", "API token")
	cmd.Flags().StringP("url", "u", "", "API base URL")
	cmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")
	viper.BindPFlags(cmd.Flags())
}

func AddCommonServeFlags(cmd *cobra.Command) {
	cmd.Flags().Int("port", 8080, "Port to serve on")
	viper.BindPFlags(cmd.Flags())
}

func AddUniversalFlags(cmd *cobra.Command) {
	cmd.Flags().String("out", "./out", "Output directory to serve files from")
	viper.BindPFlags(cmd.Flags())
}
