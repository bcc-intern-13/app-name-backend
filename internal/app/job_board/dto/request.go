package dto

type JobBoardFilter struct {
	City               string `query:"city"`
	JobField           string `query:"job_field"`
	JobType            string `query:"job_type"`
	Disability         string `query:"disability"`
	AccessibilityLabel string `query:"accessibility_label"`
	Search             string `query:"search"`
	Page               int    `query:"page"`
	Limit              int    `query:"limit"`
}
