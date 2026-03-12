package dto

type JobBoardFilter struct {
	Kota               string `query:"kota"`
	BidangKerja        string `query:"bidang_kerja"`
	TipePekerjaan      string `query:"tipe_pekerjaan"`
	Disabilitas        string `query:"disabilitas"`
	LabelAksesibilitas string `query:"label_aksesibilitas"`
	Search             string `query:"search"`
	Page               int    `query:"page"`
	Limit              int    `query:"limit"`
}
