package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

/* =======================
   Data Models
======================= */

type AnalyzeRequest struct {
	Resume         string `json:"resume"`
	JobDescription string `json:"job_description"`
}

type AnalyzeResponse struct {
	MatchScore    int      `json:"match_score"`
	MissingSkills []string `json:"missing_skills"`
	Summary       string   `json:"summary"`
}

/* =======================
   Main
======================= */

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/analyze", analyzeHandler)

	log.Printf("🚀 RoleMatch AI backend running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

/* =======================
   Handlers
======================= */

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
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	resp, err := analyze(req.Resume, req.JobDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* =======================
   Provider Switch
======================= */

func analyze(resume, jd string) (*AnalyzeResponse, error) {
	provider := os.Getenv("LLM_PROVIDER")

	switch provider {
	case "gemini":
		return analyzeWithGemini(resume, jd)
	case "ollama":
		return analyzeWithOllama(resume, jd)
	default:
		return nil, errors.New("invalid LLM_PROVIDER (use 'gemini' or 'ollama')")
	}
}

/* =======================
   Gemini Implementation
======================= */

func analyzeWithGemini(resume, jd string) (*AnalyzeResponse, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	modelName := os.Getenv("GEMINI_MODEL")

	if apiKey == "" || modelName == "" {
		return nil, errors.New("GEMINI_API_KEY or GEMINI_MODEL not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	prompt := buildPrompt(resume, jd)

	result, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	text := result.Candidates[0].Content.Parts[0].(genai.Text)

	jsonStr, err := extractJSON(string(text))
	if err != nil {
		return nil, err
	}

	var parsed AnalyzeResponse
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

/* =======================
   Ollama Implementation
======================= */

func analyzeWithOllama(resume, jd string) (*AnalyzeResponse, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_MODEL")

	if baseURL == "" || model == "" {
		return nil, errors.New("OLLAMA_BASE_URL or OLLAMA_MODEL not set")
	}

	prompt := buildPrompt(resume, jd)

	payload := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	jsonStr, err := extractJSON(result.Response)
	if err != nil {
		return nil, err
	}

	var parsed AnalyzeResponse
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

/* =======================
   Prompt Builder
======================= */

func buildPrompt(resume, jd string) string {
	return `
You are an ATS resume analyzer.

Return ONLY valid JSON in this EXACT format:
{
  "match_score": number (0-100),
  "missing_skills": [string],
  "summary": string
}

Resume:
` + resume + `

Job Description:
` + jd
}

/* =======================
   JSON Extraction Helper
======================= */

func extractJSON(text string) (string, error) {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")

	if start == -1 || end == -1 || end <= start {
		return "", errors.New("failed to extract JSON from model output")
	}

	return text[start : end+1], nil
}
