package dto

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"
	"time"
)

// =========================
// Response DTO
// =========================

type LectureSessionDTO struct {
	LectureSessionID            string    `json:"lecture_session_id"`
	LectureSessionTitle         string    `json:"lecture_session_title"`
	LectureSessionDescription   string    `json:"lecture_session_description,omitempty"`
	LectureSessionTeacherID     string    `json:"lecture_session_teacher_id"`
	LectureSessionScheduledTime time.Time `json:"lecture_session_scheduled_time"`
	LectureSessionPlace         *string   `json:"lecture_session_place,omitempty"`
	LectureSessionIsSingle      bool      `json:"lecture_session_is_single"`
	LectureSessionLectureID     *string   `json:"lecture_session_lecture_id,omitempty"`
	LectureSessionCreatedAt     time.Time `json:"lecture_session_created_at"`
}

// =========================
// Request DTOs
// =========================

type CreateLectureSessionRequest struct {
	LectureSessionTitle         string    `json:"lecture_session_title" validate:"required,min=3"`
	LectureSessionDescription   string    `json:"lecture_session_description,omitempty"`
	LectureSessionScheduledTime time.Time `json:"lecture_session_scheduled_time" validate:"required"`
	LectureSessionPlace         *string   `json:"lecture_session_place,omitempty"`
	LectureSessionIsSingle      bool      `json:"lecture_session_is_single"`
	LectureSessionLectureID     *string   `json:"lecture_session_lecture_id,omitempty"`
}

type UpdateLectureSessionRequest struct {
	LectureSessionTitle         string    `json:"lecture_session_title" validate:"required,min=3"`
	LectureSessionDescription   string    `json:"lecture_session_description,omitempty"`
	LectureSessionScheduledTime time.Time `json:"lecture_session_scheduled_time" validate:"required"`
	LectureSessionPlace         *string   `json:"lecture_session_place,omitempty"`
	LectureSessionIsSingle      bool      `json:"lecture_session_is_single"`
	LectureSessionLectureID     *string   `json:"lecture_session_lecture_id,omitempty"`
}

// =========================
// Request → Model converter
// =========================

func (r CreateLectureSessionRequest) ToModel(teacherID string) model.LectureSessionModel {
	return model.LectureSessionModel{
		LectureSessionTitle:         r.LectureSessionTitle,
		LectureSessionDescription:   r.LectureSessionDescription,
		LectureSessionScheduledTime: r.LectureSessionScheduledTime,
		LectureSessionPlace:         r.LectureSessionPlace,
		LectureSessionLectureID:     r.LectureSessionLectureID,
		LectureSessionTeacherID:     teacherID,
	}
}

// =========================
// Model → Response converter
// =========================

func ToLectureSessionDTO(m model.LectureSessionModel) LectureSessionDTO {
	return LectureSessionDTO{
		LectureSessionID:            m.LectureSessionID,
		LectureSessionTitle:         m.LectureSessionTitle,
		LectureSessionDescription:   m.LectureSessionDescription,
		LectureSessionTeacherID:     m.LectureSessionTeacherID,
		LectureSessionScheduledTime: m.LectureSessionScheduledTime,
		LectureSessionPlace:         m.LectureSessionPlace,
		LectureSessionLectureID:     m.LectureSessionLectureID,
		LectureSessionCreatedAt:     m.LectureSessionCreatedAt,
	}
}
