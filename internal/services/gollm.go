package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/teilomillet/gollm"
)

func GeneratePrompt(tagName string, commitMessages []string) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("relaise config file not created yet. Run `relaise config --api-key AI_API_KEY`.\n")
		return "", err
	}

	llm, err := gollm.NewLLM(
		gollm.SetProvider("mistral"),
		gollm.SetModel("mistral-small-latest"),
		gollm.SetAPIKey(cfg.APIKey),
		gollm.SetTemperature(0.7),
		gollm.SetMaxTokens(500),
		gollm.SetMaxRetries(3),
		gollm.SetRetryDelay(time.Second*2),
		gollm.SetLogLevel(gollm.LogLevelInfo))

	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	ctx := context.Background()

	promptText := "Generate a release note based on the following commit messages:\n\n"
	for _, msg := range commitMessages {
		promptText += "- " + msg + "\n"
	}
	promptText += fmt.Sprintf("\n Based on the amount of changes, propose next Tag Name following the release guidelines: %s\n", tagName)

	prompt := gollm.NewPrompt(promptText,
		gollm.WithDirectives("Be brief and professional", "Don't use any emojis", "Use bullet points"))

	response, err := llm.Generate(ctx, prompt)

	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}
	return response, nil
}
