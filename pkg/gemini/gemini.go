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
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiService{
		client: client,
		model:  "gemini-2.0-flash",
	}, nil
}

func (g *GeminiService) Close() {
	g.client.Close()
}

// GenerateText → prompt teks biasa, untuk improve sentence, review, dll
func (g *GeminiService) GenerateText(ctx context.Context, prompt string) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(2048)

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

	prompt := `Kamu adalah sistem ekstraksi CV profesional.

Ekstrak semua informasi dari CV PDF ini dan return HANYA dalam format JSON berikut tanpa markdown atau teks lain apapun:
{
  "summary": "ringkasan profil kandidat dalam 2-3 kalimat",
  "education": [
    {
      "institution": "nama institusi",
      "major": "jurusan",
      "year_start": "tahun masuk",
      "year_end": "tahun keluar"
    }
  ],
  "experience": [
    {
      "company": "nama perusahaan",
      "position": "posisi jabatan",
      "description": "deskripsi pekerjaan",
      "year_start": "tahun mulai",
      "year_end": "tahun selesai"
    }
  ],
  "skills": ["skill 1", "skill 2"],
  "adaptive_skills": ["skill adaptif 1", "skill adaptif 2"],
  "extracted_text": "semua teks mentah dari CV digabung jadi 1 string"
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
