package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type JobListingResponse struct {
	ID                  uuid.UUID      `json:"id"`
	CompanyID           uuid.UUID      `json:"company_id"`
	CompanyName         string         `json:"company_name"`
	CompanyLogo         string         `json:"company_logo"`
	Judul               string         `json:"judul"`
	Kota                string         `json:"kota"`
	TipePekerjaan       string         `json:"tipe_pekerjaan"`
	BidangKerja         string         `json:"bidang_kerja"`
	Gaji                string         `json:"gaji"`
	DisabilitasDiterima datatypes.JSON `json:"disabilitas_diterima"`
	LabelAksesibilitas  datatypes.JSON `json:"label_aksesibilitas"`
	CreatedAt           time.Time      `json:"created_at"`
}

type JobListingDetailResponse struct {
	JobListingResponse
	Deskripsi   string `json:"deskripsi"`
	Kualifikasi string `json:"kualifikasi"`
}

type PaginatedJobResponse struct {
	Data  []JobListingResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}
