package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, service contract.CVService, jwtSecret string) {
	h := &cvHandler{service: service}

	cv := router.Group("api/cv", middleware.JWTProtected(jwtSecret))
	cv.Post("/upload", h.uploadCV)
	cv.Post("/analyze", h.analyzeCV)
	cv.Get("/ai-calls-remaining", h.getAICallsRemaining)

	cvAI := router.Group("api/cv-ai", middleware.JWTProtected(jwtSecret))
	cvAI.Get("/score", h.getScore)
	cvAI.Post("/improve-sentence", h.ImproveSentence)
	cvAI.Post("/suggest-keywords", h.SuggestKeywords)
	cvAI.Post("/summarize-profile", h.SummarizeProfile)
}
