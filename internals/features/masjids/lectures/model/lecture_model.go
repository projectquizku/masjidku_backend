package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LectureModel struct {
	LectureID            uuid.UUID      `gorm:"column:lecture_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_id"`
	LectureTitle         string         `gorm:"column:lecture_title;type:varchar(255);not null" json:"lecture_title"`
	LectureDescription   string         `gorm:"column:lecture_description;type:text" json:"lecture_description"`
	TotalLectureSessions *int           `gorm:"column:total_lecture_sessions" json:"total_lecture_sessions,omitempty"`
	LectureStatus        bool           `gorm:"column:lecture_status;default:false" json:"lecture_status"` // false = ongoing, true = finished
	LectureCertificateID *uuid.UUID     `gorm:"column:lecture_certificate_id;type:uuid" json:"lecture_certificate_id,omitempty"`
	LectureImageURL      *string        `gorm:"column:lecture_image_url;type:text" json:"lecture_image_url,omitempty"`
	LectureTeachers      datatypes.JSON `gorm:"column:lecture_teachers;type:jsonb" json:"lecture_teachers,omitempty"`
	LectureMasjidID      uuid.UUID      `gorm:"column:lecture_masjid_id;type:uuid;not null" json:"lecture_masjid_id"`

	// Pendaftaran dan pembayaran (global untuk semua sesi)
	LectureIsRegistrationRequired bool       `gorm:"column:lecture_is_registration_required;default:false" json:"lecture_is_registration_required"`
	LectureIsPaid                 bool       `gorm:"column:lecture_is_paid;default:false" json:"lecture_is_paid"`
	LecturePrice                  *int       `gorm:"column:lecture_price" json:"lecture_price,omitempty"`
	LecturePaymentDeadline        *time.Time `gorm:"column:lecture_payment_deadline" json:"lecture_payment_deadline,omitempty"`
	LecturePaymentScope           string     `gorm:"column:lecture_payment_scope;type:varchar(10);default:'lecture'" json:"lecture_payment_scope"`

	// Umum
	LectureCapacity  int            `gorm:"column:lecture_capacity" json:"lecture_capacity"`
	LectureIsPublic  bool           `gorm:"column:lecture_is_public;default:true" json:"lecture_is_public"`
	LectureCreatedAt time.Time      `gorm:"column:lecture_created_at;autoCreateTime" json:"lecture_created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (LectureModel) TableName() string {
	return "lectures"
}
