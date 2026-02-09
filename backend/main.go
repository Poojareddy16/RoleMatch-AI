package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AnalyzeRequest struct {
	Resume         string `json:"resume"`
	JobDescription string `json:"job_description"`
}

type AnalyzeResponse struct {
	MatchScore    int      `json:"match_score"`
	MissingSkills []string `json:"missing_skills"`
	Summary       string   `json:"summary"`
}

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/analyze", analyzeHandler)

	log.Printf("🚀 Gemini backend running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := analyzeWithGemini(req.Resume, req.JobDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func analyzeWithGemini(resume, jd string) (*AnalyzeResponse, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	modelName := os.Getenv("GEMINI_MODEL")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	prompt := `
You are an ATS resume analyzer.

Return ONLY valid JSON with this structure:
{
  "match_score": number (0-100),
  "missing_skills": [string],
  "summary": string
}

Resume:
` + resume + `

Job Description:
` + jd

	result, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	text := result.Candidates[0].Content.Parts[0].(genai.Text)

	var parsed AnalyzeResponse
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}
