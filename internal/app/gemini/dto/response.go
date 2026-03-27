package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CVUploadResponse → POST /api/cv/upload
type CVUploadResponse struct {
	CvURL string `json:"cv_url"`
}

// CVResponse → GET /api/cv, POST /api/cv/analyze, PATCH /api/cv
type CVResponse struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	Summary        string          `json:"summary"`
	Education      json.RawMessage `json:"education"`
	Experience     json.RawMessage `json:"experience"`
	Skills         json.RawMessage `json:"skills"`
	AdaptiveSkills json.RawMessage `json:"adaptive_skills"`
	CvScore        int             `json:"cv_score"`
	IsAiVerified   bool            `json:"is_ai_verified"`
	AiCallsToday   int             `json:"ai_calls_today"`
	CvURL          string          `json:"cv_url"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// CVScoreResponse → GET /api/cv/score
type CVScoreResponse struct {
	Score        int    `json:"score"`
	IsAiVerified bool   `json:"is_ai_verified"`
	Label        string `json:"label"` // "Rendah" | "Sedang" | "Tinggi"
}

// AICallsRemainingResponse → GET /api/cv/ai-calls-remaining
type AICallsRemainingResponse struct {
	Remaining int `json:"remaining"`
	Used      int `json:"used"`
	Max       int `json:"max"`
}

// ImproveSentenceResponse → POST /api/cv-ai/improve-sentence
type ImproveSentenceResponse struct {
	Original     string   `json:"original"`
	Alternatives []string `json:"alternatives"`
	Remaining    int      `json:"remaining_calls"`
}

// JobMatchSection → bagian dari JobMatchResponse
type JobMatchSection struct {
	Section    string `json:"section"`
	Status     string `json:"status"` // "relevan" | "bisa_ditambah" | "kurang_relevan"
	Reasoning  string `json:"reasoning"`
	Suggestion string `json:"suggestion,omitempty"`
}

// JobMatchResponse → POST /api/cv-ai/job-match
type JobMatchResponse struct {
	MatchScore int               `json:"match_score"`
	Sections   []JobMatchSection `json:"sections"`
	Remaining  int               `json:"remaining_calls"`
}

// ReviewCVResponse → POST /api/cv-ai/review
type ReviewCVResponse struct {
	Strengths      []string `json:"strengths"`
	Improvements   []string `json:"improvements"`
	MainSuggestion string   `json:"main_suggestion"`
	Remaining      int      `json:"remaining_calls"`
}

// PerkuatKalimat
type SentenceSuggestion struct {
	Original     string `json:"original"`
	Alternative1 string `json:"alternative_1"`
	Alternative2 string `json:"alternative_2"`
}

type PerkuatKalimatResponse struct {
	Suggestions []SentenceSuggestion `json:"suggestions"`
	Remaining   int                  `json:"remaining"`
}

// SaranKeyword
type KeywordSuggestion struct {
	Keyword string `json:"keyword"`
	Alasan  string `json:"alasan"`
}

type SaranKeywordResponse struct {
	Keywords  []KeywordSuggestion `json:"keywords"`
	Remaining int                 `json:"remaining"`
}

// RingkasanProfil
type RingkasanProfilResponse struct {
	Ringkasan    string `json:"ringkasan"`
	VersiSingkat string `json:"versi_singkat"`
	Remaining    int    `json:"remaining"`
}
