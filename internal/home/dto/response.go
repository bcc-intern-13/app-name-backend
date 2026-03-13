package dto

import (
	"time"

	careerMappingDto "github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	jobdto "github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
)

type GreetingResponse struct {
	Nama      string    `json:"nama"`
	Timestamp time.Time `json:"timestamp"`
}

type HomeSummaryResponse struct {
	Greeting            GreetingResponse                        `json:"greeting"`
	RekomendasiLowongan []jobdto.JobListingResponse             `json:"rekomendasi_lowongan"`
	CareerMapping       *careerMappingDto.CareerMappingResponse `json:"career_mapping"`
}
