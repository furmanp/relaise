/*
Copyright © 2025 furmanp <przemek@furmanp.com>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"

	"github.com/furmanp/relaise/internal/services"
	"github.com/spf13/cobra"
)

var sessionConfig services.Config

var rootCmd = &cobra.Command{
	Use:   "relaise",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := services.LoadConfig()
		if err != nil || cfg == nil {
			cfg = &services.Config{}
		}

		cmd.Flags().Visit(func(f *pflag.Flag) {
			switch f.Name {
			case "mood":
				cfg.Mood = sessionConfig.Mood
			case "include-sections":
				cfg.IncludeSections = sessionConfig.IncludeSections
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

		if err != nil {
			log.Fatalf("Failed to get commit summary: %v", err)
		}

		releaseNotes, err := services.GeneratePrompt(commitSummary.TagName, commitSummary.Messages, cfg.Mood)

		if err != nil {
			log.Fatalf("Failed to generate release notes: %v", err)
		}

		fmt.Printf(releaseNotes)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&sessionConfig.Mood, "mood", "professional", "Set the tone for the release notes")
}
