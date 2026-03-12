package dto

import "gorm.io/datatypes"

type CategoryScore struct {
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Score       int            `json:"score"`
	Description string         `json:"description,omitempty"`
	FormalJobs  datatypes.JSON `json:"formal_jobs,omitempty"`
	SideJobs    datatypes.JSON `json:"side_jobs,omitempty"`
}

type CareerMappingResponse struct {
	TopCategories []CategoryScore `json:"top_categories"`
	AllScores     []CategoryScore `json:"all_scores"`
	AttemptNumber int             `json:"attempt_number"`
}
