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
	Long: `The 'config' command allows you to set up or modify the configuration used by relaise.
	This includes selecting the AI provider, model, release note tone, language, formatting style, and more.

	Configuration is saved to a YAML file in your home directory ('~/.relaise/config.yaml')
	and is automatically loaded during each run of the tool.

	You can override specific settings using flags. For example:

	relaise config --provider mistral --model mistral-medium --language en --emojis

	This ensures consistent behavior across runs without needing to specify flags every time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Configuration saved successfully to ~/.relaise/config.yaml")
		return nil
	}}

func init() {
	configCmd.Flags().StringVarP(&config.APIKey, "api-key", "k", "", "API key for LLM provider")
	configCmd.Flags().StringVarP(&config.Language, "language", "l", "en", "Language for release notes (default)")
	configCmd.Flags().StringVarP(&config.Provider, "provider", "p", "mistral", "AI Service provider (e.g., mistral, openai)")
	configCmd.Flags().StringVarP(&config.Model, "model", "M", "mistral-small-latest", "Default model to use")
	configCmd.Flags().StringVarP(&config.BulletStyle, "bullet-style", "s", "-", "Default bullet style for lists (e.g., '-', '*')")
	configCmd.Flags().StringVarP(&config.ReleaseType, "release-type", "t", "minor", "Default release type hint (minor/major/patch)")
	configCmd.Flags().StringVarP(&config.Mood, "mood", "m", "professional", "Default tone for the release notes")
	configCmd.Flags().BoolVarP(&config.IncludeSections, "include-sections", "i", false, "Default setting for including sections")
	configCmd.Flags().BoolVarP(&config.Emojis, "emojis", "e", false, "Default setting for using emojis")
	configCmd.Flags().BoolVarP(&config.Copy, "copy", "c", false, "Default setting for copying notes to clipboard")

	rootCmd.AddCommand(configCmd)
}
