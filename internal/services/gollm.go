package services

import (
	"context"
	"fmt"
	"github.com/furmanp/relaise/internal"
	"time"

	"github.com/teilomillet/gollm"
)

func getSystemPrompt() string {
	return `You are a professional release note generator.

			Your task is to generate clear, well-structured, and concise release notes 
			for a new version of a software project. You will be provided with:
			- The previous version tag (e.g., v1.2.3)
			- A list of Git commit messages since that release
			- Contextual preferences such as tone, language, bullet style, and section layout

			Your responsibilities:
			1. Analyze commit messages and identify notable changes.
			2. Determine the new semantic version based on the release type (major, minor, patch).
			3. Format the release notes according to the given preferences.
			4. Be consistent, professional, and avoid redundancy.
			5. Follow any additional instructions provided.
			6. Do not make assumptions about the project or its context beyond what is provided.

			Output only the release notes — do not include explanations or commentary.`
}

func getAdditionalConstraints(prompt internal.NotesPrompt) []string {
	var constraints []string

	if prompt.IncludeSections {
		constraints = append(constraints, "Group changes into sections: Features, Fixes, Other.")
	}

	if prompt.ReleaseType != "" {
		if prompt.TagName != "Initial Release" && prompt.TagName != "" {
			constraints = append(constraints, fmt.Sprintf("This is a %s release. Provide the appropriate next semantic version number based on the previous tag: %s.", prompt.ReleaseType, prompt.TagName))
		} else {
			constraints = append(constraints, fmt.Sprintf("This is the first release, suggested type is %s. Generate an appropriate initial semantic version (e.g., 0.1.0 for patch/minor, 1.0.0 for major) based on the provided commits.", prompt.ReleaseType))
		}
	}

	if prompt.BulletStyle != "" {
		constraints = append(constraints, fmt.Sprintf("Use '%s' for bullets in lists.", prompt.BulletStyle))
	}

	if prompt.Language != "" {
		constraints = append(constraints, fmt.Sprintf("Write the release notes in %s.", prompt.Language))
	}

	if prompt.Emojis {
		constraints = append(constraints, "Include relevant emojis next to each item to enhance readability.")
	}

	if prompt.Mood != "" {
		constraints = append(constraints, fmt.Sprintf("Set the tone of voice to '%s'.", prompt.Mood))
	}

	constraints = append(constraints, "Avoid listing trivial or repetitive commits.")
	constraints = append(constraints, "Do not include internal or build-related changes unless significant.")

	return constraints
}

func GeneratePrompt(notestPrompt internal.NotesPrompt) (string, error) {

	llm, err := gollm.NewLLM(
		gollm.SetProvider(notestPrompt.Provider),
		gollm.SetModel(notestPrompt.Model),
		gollm.SetAPIKey(notestPrompt.APIKey),
		gollm.SetTemperature(notestPrompt.Temperature),
		gollm.SetMaxTokens(notestPrompt.MaxTokens),
		gollm.SetMaxRetries(3),
		gollm.SetRetryDelay(time.Second*2),
		gollm.SetLogLevel(gollm.LogLevelError))

	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	ctx := context.Background()

	systemPrompt := getSystemPrompt()
	promptText := fmt.Sprintf("Generate a release notes, based on the following commit messages:\n\n")

	for _, msg := range notestPrompt.Context {
		promptText += "- " + msg + "\n"
	}

	prompt := gollm.NewPrompt(
		promptText,
		gollm.WithContext(systemPrompt),
		gollm.WithDirectives(getAdditionalConstraints(notestPrompt)...))

	response, err := llm.Generate(ctx, prompt)

	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}
	return response, nil
}
