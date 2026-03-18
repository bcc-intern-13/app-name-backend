package dto

import "gorm.io/datatypes"

type OnboardingResponse struct {
	Name                    string         `json:"name"`
	Age                     string         `json:"age"`
	City                    string         `json:"city"`
	Education               string         `json:"education"`
	JobField                string         `json:"job_field"`
	JobType                 string         `json:"job_type"`
	Status                  string         `json:"status"`
	CommunicationPreference string         `json:"communication_preference"`
	WorkEnvironment         datatypes.JSON `json:"work_environment"`
	SpecialNeeds            datatypes.JSON `json:"special_needs"`
}
