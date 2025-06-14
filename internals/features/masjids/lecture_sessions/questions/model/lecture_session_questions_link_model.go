package model

import (
	"time"
)

type LectureSessionsQuestionLinkModel struct {
	LectureSessionsQuestionLinkID string     `gorm:"column:lecture_sessions_question_link_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_sessions_question_link_id"`
	QuestionID                    string     `gorm:"column:question_id;type:uuid;not null" json:"question_id"`
	ExamID                        *string    `gorm:"column:exam_id;type:uuid" json:"exam_id,omitempty"`
	QuizID                        *string    `gorm:"column:quiz_id;type:uuid" json:"quiz_id,omitempty"`
	QuestionOrder                 *int       `gorm:"column:question_order" json:"question_order,omitempty"`
	CreatedAt                     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (LectureSessionsQuestionLinkModel) TableName() string {
	return "lecture_sessions_question_links"
}
