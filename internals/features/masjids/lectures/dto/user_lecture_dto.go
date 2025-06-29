package dto

import (
	"masjidku_backend/internals/features/masjids/lectures/model"
	"time"

	"github.com/google/uuid"
)

// ============================
// ✅ REQUEST STRUCT
// ============================
type UserLectureRequest struct {
	UserLectureUserID                 uuid.UUID  `json:"user_lecture_user_id"`
	UserLectureLectureID              uuid.UUID  `json:"user_lecture_lecture_id"`
	UserLectureGradeResult            *int       `json:"user_lecture_grade_result,omitempty"`
	UserLectureTotalCompletedSessions int        `json:"user_lecture_total_completed_sessions"`
	UserLectureIsRegistered           bool       `json:"user_lecture_is_registered"`
	UserLectureHasPaid                bool       `json:"user_lecture_has_paid"`
	UserLecturePaidAmount             *int       `json:"user_lecture_paid_amount,omitempty"`
	UserLecturePaymentTime            *time.Time `json:"user_lecture_payment_time,omitempty"`
}

// ============================
// ✅ RESPONSE STRUCT
// ============================
type UserLectureResponse struct {
	UserLectureID                     uuid.UUID `json:"user_lecture_id"`
	UserLectureUserID                 uuid.UUID `json:"user_lecture_user_id"`
	UserLectureLectureID              uuid.UUID `json:"user_lecture_lecture_id"`
	UserLectureGradeResult            *int      `json:"user_lecture_grade_result,omitempty"`
	UserLectureTotalCompletedSessions int       `json:"user_lecture_total_completed_sessions"`
	UserLectureIsRegistered           bool      `json:"user_lecture_is_registered"`
	UserLectureHasPaid                bool      `json:"user_lecture_has_paid"`
	UserLecturePaidAmount             *int      `json:"user_lecture_paid_amount,omitempty"`
	UserLecturePaymentTime            *string   `json:"user_lecture_payment_time,omitempty"`
	UserLectureCreatedAt              string    `json:"user_lecture_created_at"`
}

// ============================
// 🔁 Convert request → model
// ============================
func (r *UserLectureRequest) ToModel() *model.UserLectureModel {
	return &model.UserLectureModel{
		UserLectureUserID:                 r.UserLectureUserID,
		UserLectureLectureID:              r.UserLectureLectureID,
		UserLectureGradeResult:            r.UserLectureGradeResult,
		UserLectureTotalCompletedSessions: r.UserLectureTotalCompletedSessions,
		UserLectureIsRegistered:           r.UserLectureIsRegistered,
		UserLectureHasPaid:                r.UserLectureHasPaid,
		UserLecturePaidAmount:             r.UserLecturePaidAmount,
		UserLecturePaymentTime:            r.UserLecturePaymentTime,
	}
}

// ============================
// 🔁 Convert model → response
// ============================
func ToUserLectureResponse(m *model.UserLectureModel) *UserLectureResponse {
	var paymentTimeStr *string
	if m.UserLecturePaymentTime != nil {
		s := m.UserLecturePaymentTime.Format("2006-01-02 15:04:05")
		paymentTimeStr = &s
	}

	return &UserLectureResponse{
		UserLectureID:                     m.UserLectureID,
		UserLectureUserID:                 m.UserLectureUserID,
		UserLectureLectureID:              m.UserLectureLectureID,
		UserLectureGradeResult:            m.UserLectureGradeResult,
		UserLectureTotalCompletedSessions: m.UserLectureTotalCompletedSessions,
		UserLectureIsRegistered:           m.UserLectureIsRegistered,
		UserLectureHasPaid:                m.UserLectureHasPaid,
		UserLecturePaidAmount:             m.UserLecturePaidAmount,
		UserLecturePaymentTime:            paymentTimeStr,
		UserLectureCreatedAt:              m.UserLectureCreatedAt.Format("2006-01-02 15:04:05"),
	}
}
