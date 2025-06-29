package dto

import (
	"encoding/json"
	"masjidku_backend/internals/features/masjids/lectures/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Struct Teacher untuk frontend & penyimpanan JSON
type Teacher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ============================
// ✅ REQUEST STRUCT
// ============================
type LectureRequest struct {
	LectureTitle                  string     `json:"lecture_title"`
	LectureDescription            string     `json:"lecture_description"`
	LectureMasjidID               uuid.UUID  `json:"lecture_masjid_id"`
	TotalLectureSessions          *int       `json:"total_lecture_sessions"`
	LectureImageURL               *string    `json:"lecture_image_url"`
	LectureTeachers               []Teacher  `json:"lecture_teachers"`
	LectureStatus                 bool       `json:"lecture_status"` // false = ongoing, true = finished
	LectureCertificateID          *uuid.UUID `json:"lecture_certificate_id,omitempty"`
	LectureIsRegistrationRequired bool       `json:"lecture_is_registration_required"`
	LectureIsPaid                 bool       `json:"lecture_is_paid"`
	LecturePrice                  *int       `json:"lecture_price,omitempty"`
	LecturePaymentDeadline        *string    `json:"lecture_payment_deadline,omitempty"` // string format
	LecturePaymentScope           string     `json:"lecture_payment_scope"`
	LectureCapacity               int        `json:"lecture_capacity"`
	LectureIsPublic               bool       `json:"lecture_is_public"`
}

// ============================
// ✅ RESPONSE STRUCT
// ============================
type LectureResponse struct {
	LectureID                     uuid.UUID  `json:"lecture_id"`
	LectureTitle                  string     `json:"lecture_title"`
	LectureDescription            string     `json:"lecture_description"`
	LectureMasjidID               uuid.UUID  `json:"lecture_masjid_id"`
	TotalLectureSessions          *int       `json:"total_lecture_sessions"`
	LectureImageURL               *string    `json:"lecture_image_url"`
	LectureTeachers               []Teacher  `json:"lecture_teachers"`
	LectureStatus                 bool       `json:"lecture_status"`
	LectureCertificateID          *uuid.UUID `json:"lecture_certificate_id,omitempty"`
	LectureIsRegistrationRequired bool       `json:"lecture_is_registration_required"`
	LectureIsPaid                 bool       `json:"lecture_is_paid"`
	LecturePrice                  *int       `json:"lecture_price,omitempty"`
	LecturePaymentDeadline        *string    `json:"lecture_payment_deadline,omitempty"`
	LecturePaymentScope           string     `json:"lecture_payment_scope"`
	LectureCapacity               int        `json:"lecture_capacity"`
	LectureIsPublic               bool       `json:"lecture_is_public"`
	LectureCreatedAt              string     `json:"lecture_created_at"`
}

// ============================
// 🔁 CONVERT REQUEST → MODEL
// ============================
func (r *LectureRequest) ToModel() *model.LectureModel {
	teacherJSON, _ := json.Marshal(r.LectureTeachers)

	var deadline *time.Time
	if r.LecturePaymentDeadline != nil {
		t, err := time.Parse("2006-01-02 15:04:05", *r.LecturePaymentDeadline)
		if err == nil {
			deadline = &t
		}
	}

	return &model.LectureModel{
		LectureTitle:                  r.LectureTitle,
		LectureDescription:            r.LectureDescription,
		LectureMasjidID:               r.LectureMasjidID,
		TotalLectureSessions:          r.TotalLectureSessions,
		LectureImageURL:               r.LectureImageURL,
		LectureTeachers:               datatypes.JSON(teacherJSON),
		LectureStatus:                 r.LectureStatus,
		LectureCertificateID:          r.LectureCertificateID,
		LectureIsRegistrationRequired: r.LectureIsRegistrationRequired,
		LectureIsPaid:                 r.LectureIsPaid,
		LecturePrice:                  r.LecturePrice,
		LecturePaymentDeadline:        deadline,
		LecturePaymentScope:           r.LecturePaymentScope,
		LectureCapacity:               r.LectureCapacity,
		LectureIsPublic:               r.LectureIsPublic,
	}
}

// ============================
// 🔁 CONVERT MODEL → RESPONSE
// ============================
func ToLectureResponse(m *model.LectureModel) *LectureResponse {
	var teachers []Teacher
	if m.LectureTeachers != nil {
		_ = json.Unmarshal(m.LectureTeachers, &teachers)
	}

	var deadlineStr *string
	if m.LecturePaymentDeadline != nil {
		s := m.LecturePaymentDeadline.Format("2006-01-02 15:04:05")
		deadlineStr = &s
	}

	return &LectureResponse{
		LectureID:                     m.LectureID,
		LectureTitle:                  m.LectureTitle,
		LectureDescription:            m.LectureDescription,
		LectureMasjidID:               m.LectureMasjidID,
		TotalLectureSessions:          m.TotalLectureSessions,
		LectureImageURL:               m.LectureImageURL,
		LectureTeachers:               teachers,
		LectureStatus:                 m.LectureStatus,
		LectureCertificateID:          m.LectureCertificateID,
		LectureIsRegistrationRequired: m.LectureIsRegistrationRequired,
		LectureIsPaid:                 m.LectureIsPaid,
		LecturePrice:                  m.LecturePrice,
		LecturePaymentDeadline:        deadlineStr,
		LecturePaymentScope:           m.LecturePaymentScope,
		LectureCapacity:               m.LectureCapacity,
		LectureIsPublic:               m.LectureIsPublic,
		LectureCreatedAt:              m.LectureCreatedAt.Format("2006-01-02 15:04:05"),
	}
}
