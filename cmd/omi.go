package cmd

import (
	"fmt"
	"os"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
)

var omiCmd = &cobra.Command{
	Use:   "omi",
	Short: "Process Omi export text files",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, _ := cmd.Flags().GetString("source")
		outDir, _ := cmd.Flags().GetString("out")

		err := common.ProcessTextExports(sourceDir, outDir, "omi", common.ParseOmiFile)
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
