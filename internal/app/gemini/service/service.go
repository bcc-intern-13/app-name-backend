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
		repo:    repo,
		gemini:  geminiSvc,
		storage: storageSvc,
	}
}

// UploadCV → terima PDF → upload ke Supabase Storage → simpan cv_url ke cvs (TANPA Gemini)
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

	// upload ke Supabase Storage
	cvURL, err := s.storage.UploadCV(userID.String(), fileBytes, "application/pdf")
	if err != nil {
		slog.Error("failed to upload cv to storage", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to upload cv")
	}

	// upsert cv_url ke tabel cvs (belum ada extracted data)
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

// AnalyzeCV → ambil cv_url dari cvs → fetch PDF → Gemini extract → update cvs
func (s *cvService) AnalyzeCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil || cv.CvURL == "" {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	// fetch PDF dari Supabase Storage
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

	// update cv dengan hasil extract
	cv.Summary = extracted.Summary
	cv.Education = datatypes.JSON(mustMarshal(extracted.Education))
	cv.Experience = datatypes.JSON(mustMarshal(extracted.Experience))
	cv.Skills = datatypes.JSON(mustMarshal(extracted.Skills))
	cv.AdaptiveSkills = datatypes.JSON(mustMarshal(extracted.AdaptiveSkills))
	cv.CvScore = calculateScore(cv)
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

func (s *cvService) UpdateCV(ctx context.Context, userID uuid.UUID, req *dto.UpdateCVRequest) (*dto.CVResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil {
		return nil, response.ErrNotFound("cv not found, please upload your cv first")
	}

	if req.Summary != nil {
		cv.Summary = *req.Summary
	}
	if req.Education != nil {
		cv.Education = datatypes.JSON(mustMarshal(req.Education))
	}
	if req.Experience != nil {
		cv.Experience = datatypes.JSON(mustMarshal(req.Experience))
	}
	if req.Skills != nil {
		cv.Skills = datatypes.JSON(mustMarshal(req.Skills))
	}
	if req.AdaptiveSkills != nil {
		cv.AdaptiveSkills = datatypes.JSON(mustMarshal(req.AdaptiveSkills))
	}

	cv.CvScore = calculateScore(cv)
	cv.IsAiVerified = cv.CvScore >= 80

	if err := s.repo.Update(ctx, cv); err != nil {
		slog.Error("failed to update cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to update cv")
	}

	return toCVResponse(cv), nil
}

func (s *cvService) GetScore(ctx context.Context, userID uuid.UUID) (*dto.CVScoreResponse, *response.APIError) {
	cv, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get cv", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get cv")
	}
	if cv == nil {
		return nil, response.ErrNotFound("cv not found")
	}

	score := calculateScore(cv)
	cv.CvScore = score
	cv.IsAiVerified = score >= 80
	_ = s.repo.Update(ctx, cv)

	return &dto.CVScoreResponse{
		Score:        score,
		IsAiVerified: cv.IsAiVerified,
		Label:        scoreLabel(score),
	}, nil
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

// ── AI Features (Premium) ─────────────────────────────────────────────────────

// ── Helpers ───────────────────────────────────────────────────────────────────

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

func calculateScore(cv *entity.CV) int {
	score := 0
	if cv.Summary != "" {
		score += 5
	}
	if len(cv.Education) > 2 {
		score += 5
	}
	if len(cv.Experience) > 2 {
		score += 10
	}
	if len(cv.Skills) > 2 {
		score += 5
	}
	if len(cv.Summary) > 100 {
		score += 10
	}
	if containsNumber(cv.Summary) {
		score += 15
	}
	if len(cv.Skills) > 10 {
		score += 15
	}
	if len(cv.AdaptiveSkills) > 2 {
		score += 10
	}
	if len(cv.Experience) > 50 {
		score += 25
	}
	if score > 100 {
		score = 100
	}
	return score
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

func containsNumber(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
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

// PerkuatKalimat → identifikasi kalimat lemah di CV + kasih 2 alternatif
func (s *cvService) ImproveSentence(ctx context.Context, userID uuid.UUID) (*dto.PerkuatKalimatResponse, *response.APIError) {
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

	jsonStr, err := s.gemini.PerkuatKalimat(ctx, pdfBytes)
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

	return &dto.PerkuatKalimatResponse{
		Suggestions: parsed.Suggestions,
		Remaining:   maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}

// SuggestKeywords → identifikasi keyword penting yang belum ada di CV

func (s *cvService) SuggestKeywords(ctx context.Context, userID uuid.UUID) (*dto.SaranKeywordResponse, *response.APIError) {
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

	jsonStr, err := s.gemini.SaranKeyword(ctx, pdfBytes)
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

	return &dto.SaranKeywordResponse{
		Keywords:  parsed.Keywords,
		Remaining: maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}

// RingkasanProfil → generate ringkasan profil profesional dari CV
func (s *cvService) SummarizeProfile(ctx context.Context, userID uuid.UUID) (*dto.RingkasanProfilResponse, *response.APIError) {
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

	jsonStr, err := s.gemini.RingkasanProfil(ctx, pdfBytes)
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

	return &dto.RingkasanProfilResponse{
		Ringkasan:    parsed.Ringkasan,
		VersiSingkat: parsed.VersiSingkat,
		Remaining:    maxAICallsPerDay - cv.AiCallsToday,
	}, nil
}
