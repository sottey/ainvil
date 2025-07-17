package cmd

import (
	"fmt"

	"github.com/sottey/ainvil/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var limitlessCmd = &cobra.Command{
	Use:   "limitless",
	Short: "Import data from Limitless API",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("token")
		apiURL := viper.GetString("url")
		start := viper.GetString("start")
		outputDir := viper.GetString("out")

		err := common.ParseLimitlessData(apiKey, apiURL, start, outputDir)

		if err != nil {
			fmt.Println("Error handling Limitless data: ", err)
		} else {
			fmt.Println("Done.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(limitlessCmd)
	limitlessCmd.Flags().StringP("token", "t", "", "API token for Limitless")
	limitlessCmd.Flags().StringP("url", "u", "", "API base URL")
	limitlessCmd.Flags().StringP("start", "s", "", "Optional start date (YYYY-MM-DD)")
	limitlessCmd.Flags().StringP("out", "o", "out", "Output directory")
	viper.BindPFlag("token", limitlessCmd.Flags().Lookup("token"))
	viper.BindPFlag("url", limitlessCmd.Flags().Lookup("url"))
	viper.BindPFlag("start", limitlessCmd.Flags().Lookup("start"))
	viper.BindPFlag("out", limitlessCmd.Flags().Lookup("out"))
}
