package cmd

import (
	"fmt"
	"os"

	"github.com/sottey/ainvil/common"

	"github.com/spf13/cobra"
)

var beeCmd = &cobra.Command{
	Use:   "bee",
	Short: "Process Bee export text files",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, _ := cmd.Flags().GetString("source")
		outDir, _ := cmd.Flags().GetString("out")

		err := common.ProcessTextExports(sourceDir, outDir, "bee", common.ParseBeeFile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(beeCmd)
	beeCmd.Flags().String("source", "", "Directory containing Bee .txt files (required)")
	beeCmd.Flags().String("out", "./out", "Output root directory")
}
