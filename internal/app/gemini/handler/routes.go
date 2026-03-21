package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
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

	cvAI := router.Group("api/cv-ai", middleware.JWTProtected(jwtSecret))
	cvAI.Post("/improve-sentence", h.improveSentence) // POST /api/cv-ai/improve-sentence
	cvAI.Post("/job-match", h.jobMatch)               // POST /api/cv-ai/job-match
	cvAI.Post("/review", h.reviewCV)                  // POST /api/cv-ai/review
}
