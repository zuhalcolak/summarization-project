package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zuhalcolak/summarization-project/handlers"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	http.HandleFunc("/summarization", handlers.GetSummarizationText)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
