package dto

import (
	"time"

	careerMappingContract "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/dto"
	"gorm.io/datatypes"
)

type PersonalInfoResponse struct {
	Name      string `json:"name"`
	Age       string `json:"age"`
	City      string `json:"city"`
	Education string `json:"education"`
}

type CareerPreferenceResponse struct {
	JobField string `json:"job_field"`
	JobType  string `json:"job_type"`
	Status   string `json:"status"`
}

type CommunicationResponse struct {
	CommunicationPreference string `json:"communication_preference"`
}

type AccessibilityResponse struct {
	WorkEnvironment datatypes.JSON `json:"work_environment"`
	SpecialNeeds    datatypes.JSON `json:"special_needs"`
}

type SmartProfileResponse struct {
	PersonalInfo     PersonalInfoResponse                         `json:"personal_info"`
	CareerPreference CareerPreferenceResponse                     `json:"career_preference"`
	Communication    CommunicationResponse                        `json:"communication"`
	Accessibility    AccessibilityResponse                        `json:"accessibility"`
	CareerMapping    *careerMappingContract.CareerMappingResponse `json:"career_mapping"`
	UpdatedAt        time.Time                                    `json:"updated_at"`
}
