package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, service contract.CVService, jwtSecret string) {
	h := &cvHandler{service: service}

	cv := router.Group("api/cv", middleware.JWTProtected(jwtSecret))
	cv.Post("/upload", h.uploadCV)                       // POST  /api/cv/upload   → upload PDF ke storage
	cv.Post("/analyze", h.analyzeCV)                     // POST  /api/cv/analyze  → Gemini extract
	cv.Get("", h.getCV)                                  // GET   /api/cv
	cv.Patch("", h.updateCV)                             // PATCH /api/cv
	cv.Get("/score", h.getScore)                         // GET   /api/cv/score
	cv.Get("/ai-calls-remaining", h.getAICallsRemaining) // GET   /api/cv/ai-calls-remaining

	cvAI := router.Group("api/cv-ai", middleware.JWTProtected(jwtSecret)) // POST /api/cv-ai/review
	// // di dalam RegisterCVRoutes, tambah 3 route baru (semua premium + JWT protected):
	cvAI.Post("/improve-sentence", middleware.JWTProtected(jwtSecret), h.ImproveSentence)
	cvAI.Post("/suggest-keywords", middleware.JWTProtected(jwtSecret), h.SuggestKeywords)
	cvAI.Post("/summarize-profile", middleware.JWTProtected(jwtSecret), h.SummarizeProfile)
}
