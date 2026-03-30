package dto

import (
	"time"

	careerMappingDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/dto"
	jobdto "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
)

type GreetingResponse struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	AvatarURL string    `json:"avatar_url"`
}

type HomeSummaryResponse struct {
	Greeting           GreetingResponse                        `json:"greeting"`
	JobRecommendations []jobdto.JobListingResponse             `json:"job_recommendations"`
	CareerMapping      *careerMappingDto.CareerMappingResponse `json:"career_mapping"`
}
