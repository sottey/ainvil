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
