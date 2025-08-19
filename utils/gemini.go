package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultModel = "models/gemini-2.5-flash"

type GeminiClient struct {
	APIKey string
	Model  string
	Http   *http.Client
}

func NewGeminiClient() (*GeminiClient, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}
	return &GeminiClient{
		APIKey: key,
		Model:  defaultModel,
		Http:   &http.Client{Timeout: 45 * time.Second},
	}, nil
}

type part struct {
	Text string `json:"text"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type sysContent struct {
	Parts []part `json:"parts"`
}

type requestBody struct {
	SystemInstruction *sysContent      `json:"systemInstruction,omitempty"`
	Contents          []content        `json:"contents"`
	Tools             []map[string]any `json:"tools,omitempty"`
	SafetySettings    []map[string]any `json:"safetySettings,omitempty"`
	GenerationConfig  map[string]any   `json:"generationConfig,omitempty"`
}

type candidatePart struct {
	Text string `json:"text"`
}

type candidateContent struct {
	Parts []candidatePart `json:"parts"`
}

type candidate struct {
	Content candidateContent `json:"content"`
}

type geminiResponse struct {
	Candidates []candidate `json:"candidates"`
}

func (c *GeminiClient) GenerateText(systemInstruction, userPrompt string) (string, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s",
		c.Model, c.APIKey,
	)

	body := requestBody{
		SystemInstruction: &sysContent{
			Parts: []part{{Text: systemInstruction}},
		},
		Contents: []content{
			{Role: "user", Parts: []part{{Text: "userPrompt"}}},
		},
		Tools: []map[string]any{{"googleSearch": map[string]any{}}},
	}

	buf, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		trim := string(respBytes)
		if len(trim) > 300 {
			trim = trim[:300] + "..."
		}
		return "", fmt.Errorf("gemini status: %s; body: %s", resp.Status, strings.TrimSpace(trim))
	}

	var gr geminiResponse
	if err := json.Unmarshal(respBytes, &gr); err != nil {
		return "", err
	}
	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no candidates")
	}
	return strings.TrimSpace(gr.Candidates[0].Content.Parts[0].Text), nil
}
