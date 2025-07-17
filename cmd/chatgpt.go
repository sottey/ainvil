package cmd

import (
	"fmt"
	"os"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
)

var ChatGPTCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "Import and normalize ChatGPT meeting transcript files",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		out, _ := cmd.Flags().GetString("out")

		if source == "" || out == "" {
			fmt.Println("Both --source and --out are required.")
			os.Exit(1)
		}

		if err := common.ParseChatGPTTranscripts(source, out); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	ChatGPTCmd.Flags().String("source", "", "Directory of ChatGPT transcript files")
	ChatGPTCmd.Flags().String("out", "", "Directory to save normalized files")
	rootCmd.AddCommand(ChatGPTCmd)
}
