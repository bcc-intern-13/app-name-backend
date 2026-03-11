package dto

import "gorm.io/datatypes"

type OnboardingResponse struct {
	Nama                 string         `json:"nama"`
	Usia                 string         `json:"usia"`
	Kota                 string         `json:"kota"`
	Pendidikan           string         `json:"pendidikan"`
	BidangKerja          string         `json:"bidang_kerja"`
	TipePekerjaan        string         `json:"tipe_pekerjaan"`
	Status               string         `json:"status"`
	PreferensiKomunikasi string         `json:"preferensi_komunikasi"`
	LingkunganKerja      datatypes.JSON `json:"lingkungan_kerja"`
	KebutuhanKhusus      datatypes.JSON `json:"kebutuhan_khusus"`
}
