package dto

import (
	"masjidku_backend/internals/features/masjids/lectures/model"

	"github.com/google/uuid"
)

// Request dari frontend → backend
type LectureRequest struct {
	LectureTitle              string    `json:"lecture_title"`
	LectureDescription        string    `json:"lecture_description"`
	LectureMasjidID           uuid.UUID `json:"lecture_masjid_id"`
	TotalLectureSessions      *int      `json:"total_lecture_sessions"`      // nullable
	LectureIsRecurring        bool      `json:"lecture_is_recurring"`        // apakah berulang
	LectureRecurrenceInterval *int      `json:"lecture_recurrence_interval"` // nullable: jumlah hari antar pertemuan
}

// Response ke frontend
type LectureResponse struct {
	LectureID                 uuid.UUID `json:"lecture_id"`
	LectureTitle              string    `json:"lecture_title"`
	LectureDescription        string    `json:"lecture_description"`
	LectureMasjidID           uuid.UUID `json:"lecture_masjid_id"`
	TotalLectureSessions      *int      `json:"total_lecture_sessions"`
	LectureIsRecurring        bool      `json:"lecture_is_recurring"`
	LectureRecurrenceInterval *int      `json:"lecture_recurrence_interval"`
	LectureCreatedAt          string    `json:"lecture_created_at"`
}

// Convert request → model
func (r *LectureRequest) ToModel() *model.LectureModel {
	return &model.LectureModel{
		LectureTitle:              r.LectureTitle,
		LectureDescription:        r.LectureDescription,
		LectureMasjidID:           r.LectureMasjidID,
		TotalLectureSessions:      r.TotalLectureSessions,
		LectureIsRecurring:        r.LectureIsRecurring,
		LectureRecurrenceInterval: r.LectureRecurrenceInterval,
	}
}

// Convert model → response
func ToLectureResponse(m *model.LectureModel) *LectureResponse {
	return &LectureResponse{
		LectureID:                 m.LectureID,
		LectureTitle:              m.LectureTitle,
		LectureDescription:        m.LectureDescription,
		LectureMasjidID:           m.LectureMasjidID,
		TotalLectureSessions:      m.TotalLectureSessions,
		LectureIsRecurring:        m.LectureIsRecurring,
		LectureRecurrenceInterval: m.LectureRecurrenceInterval,
		LectureCreatedAt:          m.LectureCreatedAt.Format("2006-01-02 15:04:05"),
	}
}
