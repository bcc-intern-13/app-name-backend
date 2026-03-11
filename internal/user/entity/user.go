package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama                string    `gorm:"type:varchar(100)" json:"nama"`
	Email               string    `gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"`
	Password            string    `gorm:"type:text;not null"                             json:"-"`
	AvatarURL           string    `json:"avatar_url" gorm:"type:varchar(255);not null"`
	IsVerified          bool      `json:"is_verified" gorm:"type:boolean;default:false;not null"`
	IsPremium           bool      `gorm:"default:false" json:"is_premium"`
	OnboardingCompleted bool      `gorm:"default:false" json:"onboarding_completed"`

	CreatedAt time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

type UserProfile struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Usia                 string         `gorm:"type:varchar(255);not null"         json:"usia"`
	Kota                 string         `gorm:"type:varchar(255);not null"         json:"kota"`
	Pendidikan           string         `gorm:"type:varchar(255);not null"         json:"pendidikan"`
	BidangKerja          string         `gorm:"type:varchar(255);not null"         json:"bidang_kerja"`
	TipePekerjaan        string         `gorm:"type:varchar(255);not null"         json:"tipe_pekerjaan"`
	Status               string         `gorm:"type:varchar(255);not null"         json:"status"`
	PreferensiKomunikasi string         `gorm:"type:varchar(255);not null"         json:"preferensi_komunikasi"`
	LingkunganKerja      datatypes.JSON `gorm:"type:jsonb" json:"lingkungan_kerja"`
	KebutuhanKhusus      datatypes.JSON `gorm:"type:jsonb" json:"kebutuhan_khusus"`
	Nama                 string         `gorm:"type:varchar(100);not null" json:"nama"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"                                 json:"updated_at"`
}
