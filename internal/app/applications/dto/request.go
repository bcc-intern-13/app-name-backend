package dto

type SubmitApplicationRequest struct {
	JobID         string `json:"job_id" validate:"required,uuid"`
	PortfolioLink string `json:"portfolio_link"`
}
