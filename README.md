Summarization Project
This project performs automatic text summarization using the Google Gemini API. It is developed with the Golang programming language.

- The input text provided by the user is sent to the Gemini model.
- The summary returned by the model is processed in a clean JSON format.
- Different texts can be quickly summarized via the API.

# Set environment variable
export GEMINI_API_KEY="YOUR_API_KEY"
https://aistudio.google.com/app/apikey

# Run the project
go run cmd/main.go
