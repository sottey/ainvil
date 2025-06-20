package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// These will be set at build time using -ldflags
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ainvil version: %s\n", Version)
		fmt.Printf("commit: %s\n", Commit)
		fmt.Printf("built at: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
