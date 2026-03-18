package dto

import "gorm.io/datatypes"

type SubmitOnboardingRequest struct {
	Name                    string         `json:"name" validate:"required,min=2,max=100"`
	Age                     string         `json:"age" validate:"required"`
	City                    string         `json:"city" validate:"required"`
	Education               string         `json:"education" validate:"required"`
	JobField                string         `json:"job_field" validate:"required"`
	JobType                 string         `json:"job_type" validate:"required"`
	Status                  string         `json:"status" validate:"required"`
	CommunicationPreference string         `json:"communication_preference" validate:"required"`
	WorkEnvironment         datatypes.JSON `json:"work_environment" validate:"required"`
	SpecialNeeds            datatypes.JSON `json:"special_needs" validate:"required"`
}
