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
