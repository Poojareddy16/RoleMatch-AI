package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

/* =======================
   Models
======================= */

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

	log.Println("🚀 Backend running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

/* =======================
   CORS
======================= */

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

/* =======================
   Handlers
======================= */

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Resume missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	_, _ = io.ReadAll(file) // Resume text optional for now

	jd := r.FormValue("job_description")
	if jd == "" {
		http.Error(w, "Job description missing", http.StatusBadRequest)
		return
	}

	resp, err := analyzeWithOllama(jd)
	if err != nil {
		log.Println("ANALYZE ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

/* =======================
   Ollama (CORRECT)
======================= */

func analyzeWithOllama(jd string) (*AnalyzeResponse, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_MODEL")

	if baseURL == "" || model == "" {
		return nil, errors.New("Ollama env vars missing")
	}

	payload := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are an ATS resume analyzer. Return ONLY valid JSON.",
			},
			{
				"role": "user",
				"content": `
Return JSON:
{
  "match_score": 0-100,
  "missing_skills": [string],
  "summary": string
}

Job Description:
` + jd,
			},
		},
		"stream": false,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		baseURL+"/api/chat",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var parsed AnalyzeResponse
	if err := json.Unmarshal([]byte(result.Message.Content), &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}
