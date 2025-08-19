package utils

import (
	"context"
	"errors"
	"os"
)

const defaultModel = "models/gemini-2.5-flash"

func (c *models.GeminiClient) GenerateTextWithContext(ctx context.Context, summarizationSystem string, prompt string) (any, error) {
	panic("unimplemented")
}

func NewGeminiClient() (*models.GeminiClient, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}

}
