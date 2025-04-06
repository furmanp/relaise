package cmd

import (
	"fmt"
	"github.com/furmanp/relaise/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var config services.Config

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up configuration for relaise",
	Run: func(cmd *cobra.Command, args []string) {
		existing, _ := services.LoadConfig()
		if existing == nil {
			existing = &services.Config{}
		}

		cmd.Flags().Visit(func(f *pflag.Flag) {
			switch f.Name {
			case "api-key":
				existing.APIKey = config.APIKey
			case "provider":
				existing.Provider = config.Provider
			case "model":
				existing.Model = config.Model
			case "mood":
				existing.Mood = config.Mood
			case "release-type":
				existing.ReleaseType = config.ReleaseType
			case "bullet-style":
				existing.BulletStyle = config.BulletStyle
			case "include-sections":
				existing.IncludeSections = config.IncludeSections
			case "language":
				existing.Language = config.Language
			}
		})

		err := services.SaveConfig(existing)
		if err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}

		fmt.Println("Configuration saved successfully.")
	}}

func init() {
	configCmd.Flags().StringVar(&config.APIKey, "api-key", "", "API key for LLM provider")
	configCmd.Flags().StringVar(&config.Language, "language", "en", "Language you want the release notes to be generated in")
	configCmd.Flags().StringVar(&config.Provider, "provider", "mistral", "AI Service provider")
	configCmd.Flags().StringVar(&config.Model, "model", "mistral-small-latest", "Model to use")
	configCmd.Flags().StringVar(&config.BulletStyle, "bullet-style", "-", "Styling for the release notes")
	configCmd.Flags().StringVar(&config.ReleaseType, "release-type", "minor", "Type of release: minor/major/patch")
	configCmd.Flags().StringVar(&config.Mood, "mood", "professional", "Set the tone for the release notes")
	configCmd.Flags().BoolVar(&config.IncludeSections, "include-sections", false, "Sections to include in the release notes (true/false")

	rootCmd.AddCommand(configCmd)
}
