package dto

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/questions/model"
	"time"
)

// ============================
// Response DTO
// ============================
type LectureSessionsQuestionLinkDTO struct {
	LectureSessionsQuestionLinkID string    `json:"lecture_sessions_question_link_id"`
	QuestionID                    string    `json:"question_id"`
	ExamID                        *string   `json:"exam_id,omitempty"`
	QuizID                        *string   `json:"quiz_id,omitempty"`
	QuestionOrder                 *int      `json:"question_order,omitempty"`
	CreatedAt                     time.Time `json:"created_at"`
}

// ============================
// Create Request DTO
// ============================
type CreateLectureSessionsQuestionLinkRequest struct {
	QuestionID    string  `json:"question_id" validate:"required,uuid"`
	ExamID        *string `json:"exam_id,omitempty" validate:"omitempty,uuid"`
	QuizID        *string `json:"quiz_id,omitempty" validate:"omitempty,uuid"`
	QuestionOrder *int    `json:"question_order,omitempty" validate:"omitempty,min=1"`
}

// ============================
// Converter
// ============================
func ToLectureSessionsQuestionLinkDTO(m model.LectureSessionsQuestionLinkModel) LectureSessionsQuestionLinkDTO {
	return LectureSessionsQuestionLinkDTO{
		LectureSessionsQuestionLinkID: m.LectureSessionsQuestionLinkID,
		QuestionID:                    m.QuestionID,
		ExamID:                        m.ExamID,
		QuizID:                        m.QuizID,
		QuestionOrder:                 m.QuestionOrder,
		CreatedAt:                     m.CreatedAt,
	}
}
