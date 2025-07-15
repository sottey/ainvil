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
	"github.com/sottey/ainvil/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ainvil",
	Short: "Ainvil is a unified AI pendant log exporter",
	Long: `Ainvil exports lifelog entries from supported AI pendant devices
like Omi and Limitless, storing them in a normalized file structure.`,
}

// Execute runs the root command.
func Execute() error {
	cobra.OnInitialize(config.InitConfig)
	return rootCmd.Execute()
}

// init binds global flags and configuration setup.
func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to config file (default is $HOME/.ainvil.json)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().String("api", "omi", "Which API to use (e.g., omi, limitless)")
	viper.BindPFlag("api", rootCmd.PersistentFlags().Lookup("api"))

	rootCmd.PersistentFlags().String("token", "", "API key or token for the selected service")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
}
