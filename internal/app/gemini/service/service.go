package service

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/entity"
	userContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/gemini"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/storage"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const maxAICallsPerDay = 10

type cvService struct {
	repo     contract.CVRepository
	gemini   *gemini.GeminiService
	storage  *storage.StorageService
	userRepo userContract.UserRepository
}

func NewCVService(
	repo contract.CVRepository,
	geminiSvc *gemini.GeminiService,
	storageSvc *storage.StorageService,
	userRepoSvc userContract.UserRepository,

) contract.CVService {
	return &cvService{
		repo:     repo,
		gemini:   geminiSvc,
		storage:  storageSvc,
		userRepo: userRepoSvc,
	}
}

// UploadCV to supabase and save the url to the database local.
func (s *cvService) UploadCV(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.CVUploadResponse, *response.APIError) {
	if file.Header.Get("Content-Type") != "application/pdf" {
		return nil, response.ErrBadRequest("cv must be a PDF file")
	}
	if file.Size > 5*1024*1024 {
		return nil, response.ErrBadRequest("cv file size must be less than 5MB")
	}

	f, err := file.Open()
	if err != nil {
		slog.Error("failed to open cv file", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to process cv file")
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		slog.Error("failed to read cv file", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to read cv file")
	}

	// upload to Supabase Storage
	cvURL, err := s.storage.UploadCV(userID.String(), fileBytes, "application/pdf")
	if err != nil {
		slog.Error("failed to upload cv to storage", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to upload cv")
	}

	// upsert cv_url to tabel cvs
	existing, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to find existing cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to check existing cv")
	}

	if existing == nil {
		cv := &entity.CV{
			ID:     uuid.New(),
			UserID: userID,
			CvURL:  cvURL,
		}
		if err := s.repo.Create(ctx, cv); err != nil {
			slog.Error("failed to create cv record", "error", err, "userID", userID)
			return nil, response.ErrInternal("failed to save cv")
		}
	} else {
		existing.CvURL = cvURL
		if err := s.repo.Update(ctx, existing); err != nil {
			slog.Error("failed to update cv url", "error", err, "userID", userID)
			return nil, response.ErrInternal("failed to update cv")
		}
	}

	return &dto.CVUploadResponse{CvURL: cvURL}, nil
}

// AnalyzeCV get cv urls from cv table and fetch pdf byte and extract it
func (s *cvService) AnalyzeCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil || cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	// fetch PDF from Supabase Storage
	pdfBytes, err := fetchFromURL(cv.CvURL)
	if err != nil {
		slog.Error("failed to fetch pdf from storage", "error", err, "url", cv.CvURL)
		return nil, response.ErrInternal("failed to fetch cv from storage")
	}

	slog.Info("analyzing cv", "size_bytes", len(pdfBytes), "userID", userID)

	// Gemini extract
	jsonStr, err := s.gemini.ExtractCVFromPDF(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to extract cv from pdf", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to analyze cv")
	}

	extracted, err := parseExtractedCV(jsonStr)
	slog.Info("gemini raw response", "response", jsonStr)
	if err != nil {
		slog.Error("failed to parse gemini response", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to parse cv extraction result")
	}

	scoreStr, err := s.gemini.ScoreCV(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to score cv during analyze", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to score cv")
	}

	var scoreParsed struct {
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(cleanJSON(scoreStr)), &scoreParsed); err != nil {
		slog.Error("failed to parse score response during analyze", "error", err)
		return nil, response.ErrInternal("failed to parse cv score")
	}

	// update cv table with the extract result
	cv.Summary = extracted.Summary
	cv.Education = datatypes.JSON(mustMarshal(extracted.Education))
	cv.Experience = datatypes.JSON(mustMarshal(extracted.Experience))
	cv.Skills = datatypes.JSON(mustMarshal(extracted.Skills))
	cv.AdaptiveSkills = datatypes.JSON(mustMarshal(extracted.AdaptiveSkills))
	cv.CvScore = scoreParsed.Score
	cv.IsAiVerified = cv.CvScore >= 80

	if err := s.repo.Update(ctx, cv); err != nil {
		slog.Error("failed to update cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to save cv analysis")
	}

	return toCVResponse(cv), nil
}

func (s *cvService) GetCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}
	return toCVResponse(cv), nil
}

func (s *cvService) GetAICallsRemaining(ctx context.Context, userID uuid.UUID) (*dto.AICallsRemainingResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil {
		return nil, response.ErrNotFound("cv not found")
	}
	return &dto.AICallsRemainingResponse{
		Remaining: maxAICallsPerDay - cv.AiCallsToday,
		Used:      cv.AiCallsToday,
		Max:       maxAICallsPerDay,
	}, nil
}

func (s *cvService) checkAndIncrementAICall(ctx context.Context, userID uuid.UUID) (*entity.CV, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}
	if cv.AiCallsToday >= maxAICallsPerDay {
		return nil, response.ErrTooManyRequests("you have reached the maximum ai calls for today (10/day), try again tomorrow")
	}
	cv.AiCallsToday++
	if err := s.repo.Update(ctx, cv); err != nil {
		slog.Error("failed to increment ai calls", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to update ai call count")
	}

	user, err := s.userRepo.FindByID(userID.String())
	if err != nil || user == nil {
		return nil, response.ErrInternal("failed to get user")
	}
	if !user.IsPremium {
		return nil, response.ErrForbidden("this feature is for premium users only")
	}
	return cv, nil
}

func scoreLabel(score int) string {
	switch {
	case score >= 80:
		return "Tinggi"
	case score >= 50:
		return "Sedang"
	default:
		return "Rendah"
	}
}

func fetchFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type extractedCV struct {
	Summary        string `json:"summary"`
	Education      any    `json:"education"`
	Experience     any    `json:"experience"`
	Skills         any    `json:"skills"`
	AdaptiveSkills any    `json:"adaptive_skills"`
}

func parseExtractedCV(jsonStr string) (*extractedCV, error) {
	var result extractedCV
	if err := json.Unmarshal([]byte(cleanJSON(jsonStr)), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func toCVResponse(cv *entity.CV) *dto.CVResponse {
	return &dto.CVResponse{
		ID:             cv.ID,
		UserID:         cv.UserID,
		Summary:        cv.Summary,
		Education:      json.RawMessage(cv.Education),
		Experience:     json.RawMessage(cv.Experience),
		Skills:         json.RawMessage(cv.Skills),
		AdaptiveSkills: json.RawMessage(cv.AdaptiveSkills),
		CvScore:        cv.CvScore,
		IsAiVerified:   cv.IsAiVerified,
		AiCallsToday:   cv.AiCallsToday,
		CvURL:          cv.CvURL,
		UpdatedAt:      cv.UpdatedAt,
	}
}

// AI feature Preimum user needed
// Imporve sentence give suggestions to make cv sentences more impactful, based on the extracted data and Gemini's analysis
func (s *cvService) ImproveSentence(ctx context.Context, userID uuid.UUID) (*dto.ImproveSentenceResponse, *response.APIError) {
	cv, apiErr := s.checkAndIncrementAICall(ctx, userID)
	if apiErr != nil {
		return nil, apiErr
	}
	if cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	pdfBytes, err := fetchFromURL(cv.CvURL)
	if err != nil {
		slog.Error("failed to fetch pdf", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to fetch cv from storage")
	}

	jsonStr, err := s.gemini.ImproveSentence(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to perkuat kalimat", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to analyze cv sentences")
	}

	var parsed struct {
		Suggestions []dto.SentenceSuggestion `json:"suggestions"`
	}
	if err := json.Unmarshal([]byte(cleanJSON(jsonStr)), &parsed); err != nil {
		slog.Error("failed to parse perkuat kalimat response", "error", err)
		return nil, response.ErrInternal("failed to parse ai response")
	}

	return &dto.ImproveSentenceResponse{
		Suggestions: parsed.Suggestions,
		Remaining:   maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}

// SuggestKeywords identificate important keyword that is missing
func (s *cvService) SuggestKeywords(ctx context.Context, userID uuid.UUID) (*dto.SuggestKeywordResponse, *response.APIError) {
	cv, apiErr := s.checkAndIncrementAICall(ctx, userID)
	if apiErr != nil {
		return nil, apiErr
	}
	if cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	pdfBytes, err := fetchFromURL(cv.CvURL)
	if err != nil {
		slog.Error("failed to fetch pdf", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to fetch cv from storage")
	}

	jsonStr, err := s.gemini.SuggestKeyword(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to saran keyword", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to analyze cv keywords")
	}

	var parsed struct {
		Keywords []dto.KeywordSuggestion `json:"keywords"`
	}
	if err := json.Unmarshal([]byte(cleanJSON(jsonStr)), &parsed); err != nil {
		slog.Error("failed to parse saran keyword response", "error", err)
		return nil, response.ErrInternal("failed to parse ai response")
	}

	return &dto.SuggestKeywordResponse{
		Keywords:  parsed.Keywords,
		Remaining: maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}

// Summarize Profile  to create a concise summary of the user's profile based on the CV content and Gemini's analysis, which can be used for quick sharing or as an introduction in job applications. This will also be a premium feature with limited daily calls. The response will include both a detailed summary and a shorter version suitable for LinkedIn or resume introductions.
func (s *cvService) SummarizeProfile(ctx context.Context, userID uuid.UUID) (*dto.SummarizeProfileResponse, *response.APIError) {
	cv, apiErr := s.checkAndIncrementAICall(ctx, userID)
	if apiErr != nil {
		return nil, apiErr
	}
	if cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	pdfBytes, err := fetchFromURL(cv.CvURL)
	if err != nil {
		slog.Error("failed to fetch pdf", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to fetch cv from storage")
	}

	jsonStr, err := s.gemini.SummarizeProfiel(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to ringkasan profil", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to generate profile summary")
	}

	var parsed struct {
		Ringkasan    string `json:"ringkasan"`
		VersiSingkat string `json:"versi_singkat"`
	}
	if err := json.Unmarshal([]byte(cleanJSON(jsonStr)), &parsed); err != nil {
		slog.Error("failed to parse ringkasan profil response", "error", err)
		return nil, response.ErrInternal("failed to parse ai response")
	}

	return &dto.SummarizeProfileResponse{
		Ringkasan:    parsed.Ringkasan,
		VersiSingkat: parsed.VersiSingkat,
		Remaining:    maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}

func (s *cvService) GetScore(ctx context.Context, userID uuid.UUID) (*dto.CVScoreResponse, *response.APIError) {

	user, err := s.userRepo.FindByID(userID.String())
	if err != nil || user == nil {
		return nil, response.ErrInternal("failed to get user")
	}
	if !user.IsPremium {
		return nil, response.ErrForbidden("this feature is for premium users only")
	}

	cv, err := s.repo.FindByUserID(ctx, userID)

	if err != nil {
		return nil, response.ErrInternal("internal server error")
	}

	if cv == nil || cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}
	if cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	pdfBytes, err := fetchFromURL(cv.CvURL)
	if err != nil {
		slog.Error("failed to fetch pdf for scoring", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to fetch cv from storage")
	}

	jsonStr, err := s.gemini.ScoreCV(ctx, pdfBytes)
	if err != nil {
		slog.Error("failed to score cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to score cv")
	}
	slog.Info("score raw response", "response", jsonStr)

	var parsed struct {
		Score    int `json:"score"`
		Feedback struct {
			Kelengkapan     string `json:"kelengkapan"`
			KekuatanKalimat string `json:"kekuatan_kalimat"`
			RelevansKarir   string `json:"relevansi_karir"`
		} `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(cleanJSON(jsonStr)), &parsed); err != nil {
		slog.Error("failed to parse score cv response", "error", err)
		return nil, response.ErrInternal("failed to parse ai response")
	}

	// update cv score in DB
	cv.CvScore = parsed.Score
	cv.IsAiVerified = parsed.Score >= 80
	_ = s.repo.Update(ctx, cv)

	return &dto.CVScoreResponse{
		Score:        parsed.Score,
		IsAiVerified: cv.IsAiVerified,
		Label:        scoreLabel(parsed.Score),
		Feedback: dto.CVScoreFeedback{
			Kelengkapan:     parsed.Feedback.Kelengkapan,
			KekuatanKalimat: parsed.Feedback.KekuatanKalimat,
			RelevansKarir:   parsed.Feedback.RelevansKarir,
		},
		Remaining: maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}
