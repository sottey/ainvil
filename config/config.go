package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig loads configuration using Viper
func InitConfig() {
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		// Use explicitly specified config file
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for .ainvil.json in the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Could not determine home directory:", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".ainvil")
		viper.SetConfigType("json")
	}

	// Read environment variables with prefix AINVIL_
	viper.SetEnvPrefix("AINVIL")
	viper.AutomaticEnv()

	// Load config file (if it exists)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Only show an error if config was explicitly set
		if cfgFile != "" {
			fmt.Printf("Error reading config file %s: %v\n", cfgFile, err)
			os.Exit(1)
		}
	}
}
