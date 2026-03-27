package gemini

import (
	"context"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client *genai.Client
	model  string
	apiKey string
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiService{
		client: client,
		model:  "gemini-2.5-flash",
		apiKey: apiKey,
	}, nil
}

func (g *GeminiService) Close() {
	g.client.Close()
}

func (g *GeminiService) GenerateText(ctx context.Context, prompt string) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(8192)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		slog.Error("failed to generate text", "error", err)
		return "", err
	}

	return g.extractText(resp), nil
}

// ExtractCVFromPDF → kirim PDF langsung ke Gemini, return JSON string
func (g *GeminiService) ExtractCVFromPDF(ctx context.Context, fileBytes []byte) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.1) // low temperature → lebih akurat untuk extraction
	model.SetMaxOutputTokens(4096)

	prompt := `Ekstrak informasi CV ini ke JSON. Return HANYA JSON tanpa markdown:
{
  "summary": "ringkasan 2-3 kalimat",
  "education": [{"institution": "", "major": "", "year_start": "", "year_end": ""}],
  "experience": [{"company": "", "position": "", "description": "", "year_start": "", "year_end": ""}],
  "skills": ["skill1", "skill2"],
  "adaptive_skills": ["skill1"]
}`

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.Blob{
			MIMEType: "application/pdf",
			Data:     fileBytes,
		},
	)
	if err != nil {
		slog.Error("failed to extract cv from pdf", "error", err)
		return "", err
	}

	return g.extractText(resp), nil
}

// extractText → helper ambil teks dari response Gemini
func (g *GeminiService) extractText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	result := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}
	return result
}

func (g *GeminiService) GetKeyPrefix() string {
	return g.apiKey[:10]
}
