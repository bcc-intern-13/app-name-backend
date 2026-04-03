package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CVUploadResponse struct {
	CvURL string `json:"cv_url"`
}
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

type CVScoreResponse struct {
	Score        int             `json:"score"`
	IsAiVerified bool            `json:"is_ai_verified"`
	Label        string          `json:"label"`
	Feedback     CVScoreFeedback `json:"feedback"`
	Remaining    int             `json:"remaining_calls"`
}
type AICallsRemainingResponse struct {
	Remaining int `json:"remaining"`
	Used      int `json:"used"`
	Max       int `json:"max"`
}
type SentenceSuggestion struct {
	Original     string `json:"original"`
	Alternative1 string `json:"alternative_1"`
	Alternative2 string `json:"alternative_2"`
}

type ImproveSentenceResponse struct {
	Suggestions []SentenceSuggestion `json:"suggestions"`
	Remaining   int                  `json:"remaining"`
}

type KeywordSuggestion struct {
	Keyword string `json:"keyword"`
	Alasan  string `json:"alasan"`
}

type SuggestKeywordResponse struct {
	Keywords  []KeywordSuggestion `json:"keywords"`
	Remaining int                 `json:"remaining"`
}

type SummarizeProfileResponse struct {
	Ringkasan    string `json:"ringkasan"`
	VersiSingkat string `json:"versi_singkat"`
	Remaining    int    `json:"remaining"`
}

type CVScoreFeedback struct {
	Kelengkapan     string `json:"kelengkapan"`
	KekuatanKalimat string `json:"kekuatan_kalimat"`
	RelevansKarir   string `json:"relevansi_karir"`
}
