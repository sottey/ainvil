package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sottey/ainvil/internal/storage"
	"github.com/sottey/ainvil/lib/api"
	"github.com/sottey/ainvil/lib/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	startDate string
	endDate   string
	all       bool
	full      bool
)

// exportCmd defines the `ainvil export` command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export lifelogs from a supported AI pendant",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token")
		if token == "" {
			log.Fatal("API token must be provided with --token or in the config file.")
		}

		client, err := api.GetClient(viper.GetString("api"), token)
		if err != nil {
			log.Fatalf("Failed to initialize API client: %v", err)
		}

		var entries []model.Entry

		if all {
			entries, err = client.GetAllEntries()
		} else {
			if startDate == "" || endDate == "" {
				log.Fatal("Must provide both --start and --end (or use --all)")
			}
			entries, err = client.GetEntries(startDate, endDate)
		}

		if err != nil {
			log.Fatalf("Failed to get entries: %v", err)
		}

		for _, entry := range entries {
			err := storage.SaveEntry(entry, full)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save entry %s: %v\n", entry.ID, err)
			}
		}

		fmt.Printf("Export complete. %d entries processed.\n", len(entries))
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().BoolVar(&all, "all", false, "Export all entries")
	exportCmd.Flags().BoolVar(&full, "full", false, "Overwrite existing files")
	exportCmd.Flags().StringVar(&startDate, "start", "", "Start date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&endDate, "end", "", "End date (YYYY-MM-DD)")
}
