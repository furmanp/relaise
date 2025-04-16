package cmd

import (
	"fmt"
	"github.com/furmanp/relaise/internal"
	"github.com/spf13/pflag"
	"os"

	"github.com/atotto/clipboard"
	"github.com/furmanp/relaise/internal/services"
	"github.com/spf13/cobra"
)

var sessionConfig internal.Config

var rootCmd = &cobra.Command{
	Use:   "relaise",
	Short: "Generate AI-powered release notes from your Git commit history.",
	Long: `Relaise is a CLI tool that automatically generates release notes based on your Git history 
			and commit messages since the latest annotated semantic version tag (e.g. v1.2.3).

			It uses a local configuration file to define preferences such as language, bullet styles, 
			tone (mood), and whether to include emojis or structured sections.

			Relaise communicates with an AI model (via your configured provider, e.g. Mistral) 
			to convert your commit log into polished, human-friendly release notes.

			Basic usage:

				relaise

			You can customize the output using flags like:

				relaise --release-type patch --language fr --include-sections --emojis

			Relaise automatically detects your Git repository and collects commits since the last version tag.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := internal.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration (~/.relaise/config.yaml): %w\nRun 'relaise config --api-key YOUR_API_KEY'", err)
		}

		if cfg.APIKey == "" {
			return fmt.Errorf("AI API Key not found in configuration. Please set it using 'relaise config --api-key YOUR_API_KEY'")
		}

		prompt := internal.NotesPrompt{Config: *cfg}

		cmd.Flags().Visit(func(f *pflag.Flag) {
			switch f.Name {
			case "mood":
				prompt.Mood = sessionConfig.Mood
			case "include-sections":
				prompt.IncludeSections = sessionConfig.IncludeSections
			case "release-type":
				prompt.ReleaseType = sessionConfig.ReleaseType
			case "bullet-style":
				prompt.BulletStyle = sessionConfig.BulletStyle
			case "language":
				prompt.Language = sessionConfig.Language
			case "emojis":
				prompt.Emojis = sessionConfig.Emojis
			case "copy":
				prompt.Copy = sessionConfig.Copy

			}
		})

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		repo, err := services.GetGitRepository(repoPath)
		if err != nil {
			return fmt.Errorf("failed to open Git repository: %w", err)
		}

		commitSummary, err := services.GetReleasePayload(repo)

		if err != nil {
			return fmt.Errorf("failed to get commit summary: %w", err)
		}

		prompt.TagName = commitSummary.TagName
		prompt.Context = commitSummary.Messages

		if len(prompt.Context) == 0 {
			fmt.Printf("No new commits found since tag %s. No release notes to generate.\n", prompt.TagName)
			return nil
		}

		releaseNotes, err := services.GeneratePrompt(prompt)
		if err != nil {
			return fmt.Errorf("failed to generate release notes: %w", err)
		}

		fmt.Printf(releaseNotes)

		if prompt.Copy {
			err := clipboard.WriteAll(releaseNotes)
			if err != nil {
				return fmt.Errorf("failed to copy release notes to clipboard: %v", err)
			}
			fmt.Printf("\nRelease notes copied to clipboard.\n")
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true

	rootCmd.Flags().BoolP("help", "h", false, "Help")
	rootCmd.PersistentFlags().StringVarP(&sessionConfig.Mood, "mood", "m", "professional", "Set the tone for the release notes")
	rootCmd.PersistentFlags().StringVarP(&sessionConfig.BulletStyle, "bullet-style", "s", "-", "Define the bullet style for the release notes")
	rootCmd.PersistentFlags().StringVarP(&sessionConfig.ReleaseType, "release-type", "t", "minor", "Define the release type (major, minor, patch)")
	rootCmd.PersistentFlags().StringVarP(&sessionConfig.Language, "language", "l", "en", "Define the language for the release notes")
	rootCmd.PersistentFlags().BoolVarP(&sessionConfig.IncludeSections, "include-sections", "i", false, "Include structured sections (Features, Fixes, etc.)")
	rootCmd.PersistentFlags().BoolVarP(&sessionConfig.Emojis, "emojis", "e", false, "Use relevant emojis in the release notes")
	rootCmd.PersistentFlags().BoolVarP(&sessionConfig.Copy, "copy", "c", false, "Copy the generated release notes to the clipboard")

}
