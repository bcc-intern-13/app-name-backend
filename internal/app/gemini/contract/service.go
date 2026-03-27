package contract

import (
	"context"
	"mime/multipart"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type CVService interface {
	// UploadCV → terima PDF → simpan ke Supabase Storage → simpan cv_url ke cvs (TANPA Gemini)
	UploadCV(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.CVUploadResponse, *response.APIError)

	// AnalyzeCV → ambil PDF dari cv_url → Gemini extract → update cvs
	AnalyzeCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError)

	// GetCV → ambil data CV user
	GetCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError)

	// UpdateCV → update manual bagian CV (inline editing di FE)
	UpdateCV(ctx context.Context, userID uuid.UUID, req *dto.UpdateCVRequest) (*dto.CVResponse, *response.APIError)

	// GetScore → hitung ulang CV score rule-based
	GetScore(ctx context.Context, userID uuid.UUID) (*dto.CVScoreResponse, *response.APIError)

	// GetAICallsRemaining → sisa panggilan AI hari ini
	GetAICallsRemaining(ctx context.Context, userID uuid.UUID) (*dto.AICallsRemainingResponse, *response.APIError)

	// AI )features — premium only, masing-masing 1 call

	ImproveSentence(ctx context.Context, userID uuid.UUID) (*dto.PerkuatKalimatResponse, *response.APIError)
	SuggestKeywords(ctx context.Context, userID uuid.UUID) (*dto.SaranKeywordResponse, *response.APIError)
	SummarizeProfile(ctx context.Context, userID uuid.UUID) (*dto.RingkasanProfilResponse, *response.APIError)
}
