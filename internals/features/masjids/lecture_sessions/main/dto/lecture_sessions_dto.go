package dto

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"
	"time"

	"github.com/google/uuid"
)

// =========================
// Response DTO
// =========================

type LectureSessionDTO struct {
	LectureSessionID                     uuid.UUID  `json:"lecture_session_id"`
	LectureSessionTitle                  string     `json:"lecture_session_title"`
	LectureSessionDescription            string     `json:"lecture_session_description,omitempty"`
	LectureSessionTeacherID              uuid.UUID  `json:"lecture_session_teacher_id"`
	LectureSessionStartTime              time.Time  `json:"lecture_session_start_time"`
	LectureSessionEndTime                time.Time  `json:"lecture_session_end_time"`
	LectureSessionPlace                  *string    `json:"lecture_session_place,omitempty"`
	LectureSessionLectureID              *uuid.UUID `json:"lecture_session_lecture_id,omitempty"`
	LectureSessionMasjidID               uuid.UUID  `json:"lecture_session_masjid_id"`
	LectureSessionCapacity               *int       `json:"lecture_session_capacity,omitempty"`
	LectureSessionIsPublic               bool       `json:"lecture_session_is_public"`
	LectureSessionIsRegistrationRequired bool       `json:"lecture_session_is_registration_required"`
	LectureSessionIsPaid                 bool       `json:"lecture_session_is_paid"`
	LectureSessionPrice                  *int       `json:"lecture_session_price,omitempty"`
	LectureSessionPaymentDeadline        *time.Time `json:"lecture_session_payment_deadline,omitempty"`
	LectureSessionCreatedAt              time.Time  `json:"lecture_session_created_at"`
}

// =========================
// Request DTOs
// =========================

type CreateLectureSessionRequest struct {
	LectureSessionTitle                  string     `json:"lecture_session_title" validate:"required,min=3"`
	LectureSessionDescription            string     `json:"lecture_session_description,omitempty"`
	LectureSessionStartTime              time.Time  `json:"lecture_session_start_time" validate:"required"`
	LectureSessionEndTime                time.Time  `json:"lecture_session_end_time" validate:"required"`
	LectureSessionPlace                  *string    `json:"lecture_session_place,omitempty"`
	LectureSessionLectureID              *uuid.UUID `json:"lecture_session_lecture_id,omitempty"`
	LectureSessionMasjidID               uuid.UUID  `json:"lecture_session_masjid_id" validate:"required"`
	LectureSessionCapacity               *int       `json:"lecture_session_capacity,omitempty"`
	LectureSessionIsPublic               bool       `json:"lecture_session_is_public"`
	LectureSessionIsRegistrationRequired bool       `json:"lecture_session_is_registration_required"`
	LectureSessionIsPaid                 bool       `json:"lecture_session_is_paid"`
	LectureSessionPrice                  *int       `json:"lecture_session_price,omitempty"`
	LectureSessionPaymentDeadline        *time.Time `json:"lecture_session_payment_deadline,omitempty"`
}

type UpdateLectureSessionRequest struct {
	LectureSessionTitle                  string     `json:"lecture_session_title" validate:"required,min=3"`
	LectureSessionDescription            string     `json:"lecture_session_description,omitempty"`
	LectureSessionStartTime              time.Time  `json:"lecture_session_start_time" validate:"required"`
	LectureSessionEndTime                time.Time  `json:"lecture_session_end_time" validate:"required"`
	LectureSessionPlace                  *string    `json:"lecture_session_place,omitempty"`
	LectureSessionLectureID              *uuid.UUID `json:"lecture_session_lecture_id,omitempty"`
	LectureSessionMasjidID               uuid.UUID  `json:"lecture_session_masjid_id" validate:"required"`
	LectureSessionCapacity               *int       `json:"lecture_session_capacity,omitempty"`
	LectureSessionIsPublic               bool       `json:"lecture_session_is_public"`
	LectureSessionIsRegistrationRequired bool       `json:"lecture_session_is_registration_required"`
	LectureSessionIsPaid                 bool       `json:"lecture_session_is_paid"`
	LectureSessionPrice                  *int       `json:"lecture_session_price,omitempty"`
	LectureSessionPaymentDeadline        *time.Time `json:"lecture_session_payment_deadline,omitempty"`
}

// =========================
// Request → Model converter
// =========================

func (r CreateLectureSessionRequest) ToModel(teacherID uuid.UUID) model.LectureSessionModel {
	return model.LectureSessionModel{
		LectureSessionTitle:                  r.LectureSessionTitle,
		LectureSessionDescription:            r.LectureSessionDescription,
		LectureSessionStartTime:              r.LectureSessionStartTime,
		LectureSessionEndTime:                r.LectureSessionEndTime,
		LectureSessionPlace:                  r.LectureSessionPlace,
		LectureSessionLectureID:              r.LectureSessionLectureID,
		LectureSessionMasjidID:               r.LectureSessionMasjidID,
		LectureSessionCapacity:               r.LectureSessionCapacity,
		LectureSessionIsPublic:               r.LectureSessionIsPublic,
		LectureSessionIsRegistrationRequired: r.LectureSessionIsRegistrationRequired,
		LectureSessionIsPaid:                 r.LectureSessionIsPaid,
		LectureSessionPrice:                  r.LectureSessionPrice,
		LectureSessionPaymentDeadline:        r.LectureSessionPaymentDeadline,
		LectureSessionTeacherID:              teacherID,
	}
}

// =========================
// Model → Response converter
// =========================

func ToLectureSessionDTO(m model.LectureSessionModel) LectureSessionDTO {
	return LectureSessionDTO{
		LectureSessionID:                     m.LectureSessionID,
		LectureSessionTitle:                  m.LectureSessionTitle,
		LectureSessionDescription:            m.LectureSessionDescription,
		LectureSessionTeacherID:              m.LectureSessionTeacherID,
		LectureSessionStartTime:              m.LectureSessionStartTime,
		LectureSessionEndTime:                m.LectureSessionEndTime,
		LectureSessionPlace:                  m.LectureSessionPlace,
		LectureSessionLectureID:              m.LectureSessionLectureID,
		LectureSessionMasjidID:               m.LectureSessionMasjidID,
		LectureSessionCapacity:               m.LectureSessionCapacity,
		LectureSessionIsPublic:               m.LectureSessionIsPublic,
		LectureSessionIsRegistrationRequired: m.LectureSessionIsRegistrationRequired,
		LectureSessionIsPaid:                 m.LectureSessionIsPaid,
		LectureSessionPrice:                  m.LectureSessionPrice,
		LectureSessionPaymentDeadline:        m.LectureSessionPaymentDeadline,
		LectureSessionCreatedAt:              m.LectureSessionCreatedAt,
	}
}
