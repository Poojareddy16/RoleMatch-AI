# 🎯🤖 RoleMatch-AI
An intelligent job–resume matching platform built with Go that evaluates candidate resumes against job descriptions using explainable rule-based logic, with planned integration of Gemini API and local Ollama LLMs for semantic reasoning and hybrid AI matching.

---

## 🚀 Project Overview

RoleMatch-AI is a backend system designed to analyze how well a resume aligns with a job description. The current implementation uses deterministic, transparent algorithms to extract skills, compute match scores, and identify missing requirements. The platform is architected to evolve into a hybrid AI system by integrating Gemini API and Ollama-based local large language models.

The project mirrors real-world Applicant Tracking System (ATS) workflows while maintaining explainability and extensibility.

---

## 🧠 Current Capabilities

- Resume and job description parsing
- Rule-based skill extraction from unstructured text
- Skill normalization and deduplication
- Match score computation using heuristic logic
- Identification of matched and missing skills
- RESTful API implemented in Go
- Modular, layered backend architecture

---

## 🏗️ System Architecture
```bash
Client
  → Go REST API
    → Handler Layer
      → Service Layer
        → Rule-Based Matching Engine
          → Skill Extraction & Scoring Logic
```
The system is stateless and designed to support pluggable matching strategies.

---

## ⚙️ Algorithms & Logic Used (Current)

### 1. Text Preprocessing
- Case normalization
- Tokenization
- Whitespace and punctuation cleanup

Type: Deterministic text normalization

---

### 2. Skill Extraction
- Predefined technical skill dictionary
- Exact and partial keyword matching
- Deduplication using sets

Type: Dictionary-based skill extraction (NER-lite)

---

### 3. Skill Matching
- Resume skills and job skills represented as sets
- Intersection → matched skills
- Difference → missing skills

Type: Set-theoretic matching (O(n))

---

### 4. Match Scoring
- Heuristic score based on coverage of required skills

Conceptual formula:
(matched_skills / total_required_skills) × 100

Type: Explainable heuristic scoring

---

## 📡 API Endpoints (Current)

POST /match

Input:
- Resume text
- Job description text

Output:
- Match score
- Matched skills
- Missing skills

---

## 🧪 Example Output

- Match Score: 67%
- Matched Skills: SQL, Python, Machine Learning
- Missing Skills: Docker, Kubernetes

This output provides transparent insight into hiring relevance.

---

## 🔮 Future Enhancements

### 1. Gemini API Integration
- Semantic skill extraction beyond keywords
- Context-aware resume–JD matching
- Natural language reasoning for match explanations
- Skill inference and normalization

---

### 2. Hybrid Matching Engine
- Combine rule-based logic with LLM-driven reasoning
- Deterministic fallback for reliability
- Cost-aware and latency-aware routing

---

### 3. Explainable AI Output
- Human-readable justification for match scores
- Skill gap explanations
- Resume improvement suggestions

---

### 4. Production Readiness
- Environment-based configuration
- Model selection via config flags
- Logging and error handling
- Dockerized deployment

---

## 🛠️ Tech Stack

- Language: Go
- Backend: RESTful API
- Matching Logic: Rule-based algorithms
- Planned AI Models:
  - Gemini API
  - Ollama (local LLMs)

---

## 📌 Resume Summary (Post-Enhancement)

Built a hybrid AI-powered job–resume matching platform in Go, combining explainable rule-based skill matching with Gemini API and Ollama-based semantic reasoning to compute relevance scores and identify skill gaps.

---

## 📈 Project Status

- Current: Rule-based matching engine with integrated Ollama local LLM for reasoning and explanations (complete)
- Next: Google Gemini API integration for cloud-based semantic matching (planned)


---

## 📄 License

MIT License


