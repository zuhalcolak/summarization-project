package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zuhalcolak/summarization-project/models"
	"github.com/zuhalcolak/summarization-project/utils"
)

const summararizationSystem = `
You are an expert text summarizer. Given an input text, generate a concise summary in Turkish. Return ONLY a single valid JSON object with exactly one key: "summary".

Do not include code fences, explanations, or any extra keys.

If the input text is empty or cannot be summarized, return: {"summary": "Özet çıkarılamadı"}.

Otherwise, return a clear and concise summary in Turkish.

Examples:
{"summary": "This text explains the process of text summarization in the project."}
{"summary": "Özet çıkarılamadı"}
`

func buildSummarizationPrompt(req models.SummarizationRequest) string {
	return fmt.Sprintf(`
	You are an expert text summarizer. Given an input text, generate a concise summary in Turkish. 
Return ONLY a single valid JSON object with exactly one key: "geminiText".

Do not include code fences, explanations, or any extra keys.

If the input text is empty or cannot be summarized, return: {"geminiText": "Özet çıkarılamadı"}.

Otherwise, return a clear and concise summary in Turkish.

Examples:
{"geminiText": "Bu metin, projedeki özet çıkarma sürecini açıklar."}
{"geminiText": "Özet çıkarılamadı"}`,
		req.UserText)
}

func GetSummarizationText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SummarizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	client, err := utils.NewGeminiClient()
	if err != nil {
		http.Error(w, "Gemini init error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prompt := buildSummarizationPrompt(req)
	raw, err := client.GenerateText(summararizationSystem, prompt)
	if err != nil {
		http.Error(w, "Gemini API error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var out models.SummarizationResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &out); err != nil || out.GeminiText == "" {
		clean := strings.TrimPrefix(raw, "```json")
		clean = strings.TrimPrefix(strings.TrimSpace(clean), "```")
		clean = strings.TrimSuffix(strings.TrimSpace(clean), "```")
		if err2 := json.Unmarshal([]byte(strings.TrimSpace(clean)), &out); err2 != nil || out.GeminiText == "" {
			http.Error(w, "Invalid JSON from model", http.StatusBadGateway)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
