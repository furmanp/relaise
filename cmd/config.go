package cmd

import (
	"fmt"
	"github.com/furmanp/relaise/internal/services"
	"github.com/spf13/cobra"
)

var apiKey string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up configuration for relaise",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &services.Config{
			APIKey: apiKey,
		}
		err := services.SaveConfig(cfg)
		if err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}
		fmt.Println("Configuration saved successfully.")
	},
}

func init() {
	configCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for LLM provider")
	rootCmd.AddCommand(configCmd)
}
