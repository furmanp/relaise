/*
Copyright © 2025 furmanp <przemek@furmanp.com>
*/
package cmd

import (
	"fmt"
	"github.com/furmanp/relaise/internal"
	"github.com/spf13/pflag"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := internal.LoadConfig()
		if err != nil {
			log.Fatalf("AI API Key not provided. Run `relaise config --api-key AI_API_KEY`.\n")
		}

		prompt := internal.MapConfigToPrompt(*cfg)

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
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		repo, err := services.GetGitRepository(repoPath)
		if err != nil {
			log.Fatalf("Failed to open Git repository: %v", err)
		}

		commitSummary, err := services.GetReleasePayload(repo)
		prompt.TagName = commitSummary.TagName
		prompt.Context = commitSummary.Messages

		if err != nil {
			log.Fatalf("Failed to get commit summary: %v", err)
		}

		releaseNotes, err := services.GeneratePrompt(prompt)

		if err != nil {
			log.Fatalf("Failed to generate release notes: %v", err)
		}

		fmt.Printf(releaseNotes)
		if prompt.Copy {
			err := clipboard.WriteAll(releaseNotes)
			if err != nil {
				log.Fatalf("Failed to copy release notes to clipboard: %v", err)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help")
	rootCmd.PersistentFlags().StringVar(&sessionConfig.Mood, "mood", "professional", "Set the tone for the release notes")
	rootCmd.PersistentFlags().StringVar(&sessionConfig.BulletStyle, "style", "*", "Define the bullet style for the release notes")
	rootCmd.PersistentFlags().StringVar(&sessionConfig.ReleaseType, "release-type", "minor", "Define the release type")
	rootCmd.PersistentFlags().StringVar(&sessionConfig.Language, "language", "en", "Define the language for the release notes")
	rootCmd.PersistentFlags().BoolVar(&sessionConfig.IncludeSections, "include-sections", false, "Include sections in the release notes")
	rootCmd.PersistentFlags().BoolVar(&sessionConfig.Emojis, "emojis", false, "Use emojis in the release notes")
	rootCmd.PersistentFlags().BoolVar(&sessionConfig.Copy, "copy", false, "Copy the release notes to clipboard")

}
