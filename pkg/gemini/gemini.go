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

func (g *GeminiService) Generate(ctx context.Context, prompt string) (string, error) {
	model := g.client.GenerativeModel(g.model)

	// setting biar response lebih konsisten
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(2048)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	result := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	return result, nil
}

func (g *GeminiService) Close() {
	g.client.Close()
}
