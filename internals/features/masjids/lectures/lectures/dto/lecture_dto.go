package dto

import (
	"masjidku_backend/internals/features/masjids/lectures/lectures/model"

	"github.com/google/uuid"
)

type LectureRequest struct {
	LectureTitle       string    `json:"lecture_title"`
	LectureDescription string    `json:"lecture_description"`
	LecturePlace       string    `json:"lecture_place"`
	LectureMasjidID    uuid.UUID `json:"lecture_masjid_id"`
}

type LectureResponse struct {
	LectureID          uuid.UUID `json:"lecture_id"`
	LectureTitle       string    `json:"lecture_title"`
	LectureDescription string    `json:"lecture_description"`
	LecturePlace       string    `json:"lecture_place"`
	LectureMasjidID    uuid.UUID `json:"lecture_masjid_id"`
	LectureCreatedAt   string    `json:"lecture_created_at"`
}

// Convert request → model
func (r *LectureRequest) ToModel() *model.LectureModel {
	return &model.LectureModel{
		LectureTitle:       r.LectureTitle,
		LectureDescription: r.LectureDescription,
		LecturePlace:       r.LecturePlace,
		LectureMasjidID:    r.LectureMasjidID,
	}
}

// Convert model → response
func ToLectureResponse(m *model.LectureModel) *LectureResponse {
	return &LectureResponse{
		LectureID:          m.LectureID,
		LectureTitle:       m.LectureTitle,
		LectureDescription: m.LectureDescription,
		LecturePlace:       m.LecturePlace,
		LectureMasjidID:    m.LectureMasjidID,
		LectureCreatedAt:   m.LectureCreatedAt.Format("2006-01-02 15:04:05"),
	}
}
