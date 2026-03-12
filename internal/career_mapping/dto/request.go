package dto

type SubmitCareerMappingRequest struct {
	Answers []string `json:"answers" validate:"required,len=20,dive,oneof=A B C D"`
}
