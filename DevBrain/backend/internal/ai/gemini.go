package ai

import (
	"context"
	"os"

	"google.golang.org/genai"
)

func NewGeminiClient() (*genai.Client, error) {

	return genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		},
	)
}