package model

import (
	"time"
)

type LectureSessionsQuestionModel struct {
	LectureSessionsQuestionID              string     `gorm:"column:lecture_sessions_question_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	LectureSessionsQuestion                string     `gorm:"column:lecture_sessions_question;type:text;not null"`
	LectureSessionsQuestionAnswer          string     `gorm:"column:lecture_sessions_question_answer;type:text;not null"`
	LectureSessionsQuestionCorrect         string     `gorm:"column:lecture_sessions_question_correct;type:char(1);not null"` // A/B/C/D
	LectureSessionsQuestionExplanation     string     `gorm:"column:lecture_sessions_question_explanation;type:text"`
	LectureSessionsQuestionLectureSessionID *string    `gorm:"column:lecture_sessions_question_lecture_session_id;type:uuid"` // nullable
	LectureSessionsQuestionExamID           *string    `gorm:"column:lecture_sessions_question_exam_id;type:uuid"`             // nullable
	LectureSessionsQuestionCreatedAt        time.Time  `gorm:"column:lecture_sessions_question_created_at;autoCreateTime"`
}

func (LectureSessionsQuestionModel) TableName() string {
	return "lecture_sessions_questions"
}
