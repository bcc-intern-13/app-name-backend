package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type JobListing struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CompanyID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	Judul               string         `gorm:"type:varchar(200);not null" json:"judul"`
	Deskripsi           string         `gorm:"type:text;not null" json:"deskripsi"`
	Kualifikasi         string         `gorm:"type:text;not null" json:"kualifikasi"`
	Kota                string         `gorm:"type:varchar(100);not null" json:"kota"`
	TipePekerjaan       string         `gorm:"type:varchar(20);not null" json:"tipe_pekerjaan"`
	BidangKerja         string         `gorm:"type:varchar(50);not null" json:"bidang_kerja"`
	Gaji                string         `gorm:"type:varchar(100)" json:"gaji"`
	DisabilitasDiterima datatypes.JSON `gorm:"type:jsonb;not null" json:"disabilitas_diterima"`
	LabelAksesibilitas  datatypes.JSON `gorm:"type:jsonb;not null" json:"label_aksesibilitas"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

type Company struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama                string         `gorm:"type:varchar(200);not null" json:"nama"`
	LogoURL             string         `gorm:"type:text" json:"logo_url"`
	Deskripsi           string         `gorm:"type:text" json:"deskripsi"`
	Industri            string         `gorm:"type:varchar(100)" json:"industri"`
	Ukuran              string         `gorm:"type:varchar(50)" json:"ukuran"`
	Lokasi              string         `gorm:"type:varchar(100)" json:"lokasi"`
	Website             string         `gorm:"type:text" json:"website"`
	DisabilitasDiterima datatypes.JSON `gorm:"type:jsonb;not null" json:"disabilitas_diterima"`
	LabelAksesibilitas  datatypes.JSON `gorm:"type:jsonb;not null" json:"label_aksesibilitas"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
