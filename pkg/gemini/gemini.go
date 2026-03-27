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

// PerkuatKalimat → identifikasi kalimat lemah + kasih 2 alternatif
func (g *GeminiService) PerkuatKalimat(ctx context.Context, fileBytes []byte) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(4096)

	prompt := `Kamu adalah AI career advisor yang membantu penyandang disabilitas
di Indonesia meningkatkan kualitas CV mereka.

Dari CV berikut, identifikasi kalimat-kalimat yang lemah, pasif,
atau kurang impactful.

Untuk setiap kalimat yang ditemukan, berikan 2 alternatif kalimat
yang lebih kuat, action-oriented, dan profesional.

Berikan output HANYA dalam format JSON berikut, tanpa teks tambahan:
{
  "suggestions": [
    {
      "original": "<kalimat asli dari CV>",
      "alternative_1": "<alternatif kalimat 1>",
      "alternative_2": "<alternatif kalimat 2>"
    }
  ]
}
Maksimal 5 kalimat. Prioritaskan kalimat yang paling perlu diperbaiki.`

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.Blob{MIMEType: "application/pdf", Data: fileBytes},
	)
	if err != nil {
		slog.Error("failed to perkuat kalimat", "error", err)
		return "", err
	}

	return g.extractText(resp), nil
}

// SaranKeyword → identifikasi keyword yang belum ada di CV
func (g *GeminiService) SaranKeyword(ctx context.Context, fileBytes []byte) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(2048)

	prompt := `Kamu adalah AI career advisor yang membantu penyandang disabilitas
di Indonesia meningkatkan kualitas CV mereka.

Dari CV berikut, identifikasi keyword penting yang belum ada atau
kurang ditonjolkan, yang sebaiknya ditambahkan agar CV lebih mudah
ditemukan oleh HRD dan sistem ATS (Applicant Tracking System).

Berikan output HANYA dalam format JSON berikut, tanpa teks tambahan:
{
  "keywords": [
    {
      "keyword": "<keyword yang disarankan>",
      "alasan": "<alasan singkat 1 kalimat mengapa keyword ini penting>"
    }
  ]
}
Berikan maksimal 8 keyword. Prioritaskan keyword paling relevan dengan bidang karir dari CV.`

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.Blob{MIMEType: "application/pdf", Data: fileBytes},
	)
	if err != nil {
		slog.Error("failed to saran keyword", "error", err)
		return "", err
	}

	return g.extractText(resp), nil
}

// RingkasanProfil → buat ringkasan profil profesional dari CV
func (g *GeminiService) RingkasanProfil(ctx context.Context, fileBytes []byte) (string, error) {
	model := g.client.GenerativeModel(g.model)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(2048)

	prompt := `Kamu adalah AI career advisor yang membantu penyandang disabilitas
di Indonesia meningkatkan kualitas CV mereka.

Dari CV berikut, buatkan ringkasan profil profesional yang singkat,
kuat, dan menarik untuk dicantumkan di bagian atas CV.

Berikan output HANYA dalam format JSON berikut, tanpa teks tambahan:
{
  "ringkasan": "<teks ringkasan profil 3-4 kalimat>",
  "versi_singkat": "<teks ringkasan profil 1-2 kalimat>"
}
Gunakan bahasa Indonesia yang profesional dan percaya diri. Hindari kalimat pasif.`

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.Blob{MIMEType: "application/pdf", Data: fileBytes},
	)
	if err != nil {
		slog.Error("failed to ringkasan profil", "error", err)
		return "", err
	}

	return g.extractText(resp), nil
}
