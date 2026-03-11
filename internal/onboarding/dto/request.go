package dto

import "gorm.io/datatypes"

type SubmitOnboardingRequest struct {
	Nama                 string         `json:"nama" validate:"required,min=2,max=100"`
	Usia                 string         `json:"usia" validate:"required"`
	Kota                 string         `json:"kota" validate:"required"`
	Pendidikan           string         `json:"pendidikan" validate:"required"`
	BidangKerja          string         `json:"bidang_kerja" validate:"required"`
	TipePekerjaan        string         `json:"tipe_pekerjaan" validate:"required"`
	Status               string         `json:"status" validate:"required"`
	PreferensiKomunikasi string         `json:"preferensi_komunikasi" validate:"required"`
	LingkunganKerja      datatypes.JSON `json:"lingkungan_kerja" validate:"required"`
	KebutuhanKhusus      datatypes.JSON `json:"kebutuhan_khusus" validate:"required"`
}
