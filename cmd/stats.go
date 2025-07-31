package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var statsOutDir string

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show total lifelog count and most recent capture date for each pendant",
	RunE: func(cmd *cobra.Command, args []string) error {
		type stat struct {
			count      int
			latestDate time.Time
		}

		stats := make(map[string]*stat)

		err := filepath.Walk(statsOutDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}

			// Extract pendant and date from path
			parts := strings.Split(path, string(filepath.Separator))
			if len(parts) < 4 {
				return nil
			}

			year := parts[len(parts)-4]
			month := parts[len(parts)-3]
			day := parts[len(parts)-2]

			dateStr := fmt.Sprintf("%s-%s-%s", year, month, day)
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return nil
			}

			fileName := info.Name()
			pendant := strings.SplitN(fileName, "_", 2)[0]

			if _, ok := stats[pendant]; !ok {
				stats[pendant] = &stat{}
			}
			stats[pendant].count++
			if date.After(stats[pendant].latestDate) {
				stats[pendant].latestDate = date
			}
			return nil
		})

		if err != nil {
			return err
		}

		// Pretty output
		fmt.Printf("%-10s | %-5s | %-15s\n", "Pendant", "Count", "Last Capture")
		fmt.Println(strings.Repeat("-", 36))

		pendants := make([]string, 0, len(stats))
		for p := range stats {
			pendants = append(pendants, p)
		}
		sort.Strings(pendants)

		for _, p := range pendants {
			s := stats[p]
			fmt.Printf("%-10s | %-5d | %s\n", p, s.count, s.latestDate.Format("2006-01-02"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringVar(&statsOutDir, "out", "./out", "Root directory of exported lifelogs")
}
