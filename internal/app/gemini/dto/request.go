package dto

type UpdateCVRequest struct {
	Summary        *string `json:"summary"`
	Education      *any    `json:"education"`
	Experience     *any    `json:"experience"`
	Skills         *any    `json:"skills"`
	AdaptiveSkills *any    `json:"adaptive_skills"`
}

type ImproveSentenceRequest struct {
	Sentence string `json:"sentence" validate:"required,min=5"`
}

type JobMatchRequest struct {
	JobDescription string `json:"job_description" validate:"required,min=10"`
}
