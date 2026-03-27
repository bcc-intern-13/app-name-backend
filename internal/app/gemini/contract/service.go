package contract

import (
	"context"
	"mime/multipart"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type CVService interface {
	// UploadCV to receive PDF and save to Supabase Storage and save cv_url to cvs (WITHOUT Gemini to upload the cv )
	UploadCV(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.CVUploadResponse, *response.APIError)

	// AnalyzeCV get  PDF from cv_url and Gemini extract to update cvs part of the data.
	AnalyzeCV(ctx context.Context, userID uuid.UUID) (*dto.CVResponse, *response.APIError)

	// GetAICallsRemaining
	// 10 calls per day
	// Reset every day.
	GetAICallsRemaining(ctx context.Context, userID uuid.UUID) (*dto.AICallsRemainingResponse, *response.APIError)

	// AI Features
	// Premium only
	// 1 Api Call per usage
	ImproveSentence(ctx context.Context, userID uuid.UUID) (*dto.ImproveSentenceResponse, *response.APIError)
	SuggestKeywords(ctx context.Context, userID uuid.UUID) (*dto.SuggestKeywordResponse, *response.APIError)
	SummarizeProfile(ctx context.Context, userID uuid.UUID) (*dto.SummarizeProfileResponse, *response.APIError)
	GetScore(ctx context.Context, userID uuid.UUID) (*dto.CVScoreResponse, *response.APIError)
}
