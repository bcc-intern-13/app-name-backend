package dto

import "gorm.io/datatypes"

type CategoryScore struct {
	Rank        int            `json:"rank,omitempty"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Score       int            `json:"score"`
	Description string         `json:"description,omitempty"`
	FormalJobs  datatypes.JSON `json:"formal_jobs,omitempty"`
	SideJobs    datatypes.JSON `json:"side_jobs,omitempty"`
}

type CareerMappingResponse struct {
	Rank          int             `json:"rank,omitempty"`
	TopCategories []CategoryScore `json:"top_categories"`
	AllScores     []CategoryScore `json:"all_scores"`
	AttemptNumber int             `json:"attempt_number"`
}
