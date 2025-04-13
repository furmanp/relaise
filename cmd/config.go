package cmd

import (
	"fmt"
	"github.com/furmanp/relaise/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var config internal.Config

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up or update the configuration for relaise.",
	Long: `The 'config' command allows you to set up or modify the configuration used by Relaise.
	This includes selecting the AI provider, model, release note tone, language, formatting style, and more.

	Configuration is saved to a YAML file in your home directory ('~/.relaise/config.yaml')
	and is automatically loaded during each run of the tool.

	You can override specific settings using flags. For example:

	relaise config --provider mistral --model mistral-medium --language en --emojis

	This ensures consistent behavior across runs without needing to specify flags every time.`,
	Run: func(cmd *cobra.Command, args []string) {
		existing, _ := internal.LoadConfig()
		if existing == nil {
			existing = internal.DefaultConfig()
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
			case "emojis":
				existing.Emojis = config.Emojis
			case "copy":
				existing.Copy = config.Copy
			}
		})

		err := internal.SaveConfig(existing)
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
	configCmd.Flags().BoolVar(&config.Emojis, "use-emojis", false, "Use emojis in the release notes (true/false)")
	configCmd.Flags().BoolVar(&config.Copy, "copy", false, "Copy the release notes to clipboard (true/false)")

	rootCmd.AddCommand(configCmd)
}
