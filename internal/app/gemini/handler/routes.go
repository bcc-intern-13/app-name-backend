package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	userContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, service contract.CVService, jwtSecret string, userRepo userContract.UserRepository) {
	h := &cvHandler{service: service}

	// ai pre usage routes
	cv := router.Group("api/v1/cv", middleware.JWTProtected(jwtSecret))
	cv.Post("/upload", h.uploadCV)
	cv.Post("/analyze", h.analyzeCV)
	cv.Get("/ai-calls-remaining", h.getAICallsRemaining)

	//cv ai routes
	cvAI := router.Group("/api/v1/cv-ai",
		middleware.JWTProtected(jwtSecret),
		middleware.PremiumRequired(userRepo),
	)
	cvAI.Get("/score", h.getScore)
	cvAI.Post("/improve-sentence", h.ImproveSentence)
	cvAI.Post("/suggest-keywords", h.SuggestKeywords)
	cvAI.Post("/summarize-profile", h.SummarizeProfile)
}
