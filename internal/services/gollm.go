package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/teilomillet/gollm"
)

func getSystemPrompt() string {
	return `You are a release note generator. You will be given a list of commit messages and you need to generate a release note based on them.
You should follow the release guidelines and be brief and professional.
You may, or might not be given additional instructions. If you are, you should follow them.`
}

func GeneratePrompt(tagName string, commitMessages []string, mood string) (string, error) {
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

	systemPrompt := getSystemPrompt()
	promptText := fmt.Sprintf("Generate a release notes for %s, based on the following commit messages:\n\n", tagName)

	for _, msg := range commitMessages {
		promptText += "- " + msg + "\n"
	}

	prompt := gollm.NewPrompt(
		promptText,
		gollm.WithContext(systemPrompt),
		gollm.WithDirectives("Be brief and professional", "Don't use any emojis", "Use bullet points", fmt.Sprintf("make the response %s", mood)))

	response, err := llm.Generate(ctx, prompt)

	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}
	return response, nil
}
