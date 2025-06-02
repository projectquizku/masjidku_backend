package model

import (
	"time"

	"github.com/google/uuid"
)

type UserAnswerAttempt struct {
	UserAnswerAttemptID              uint      `gorm:"primaryKey" json:"user_answer_attempt_id"`
	UserAnswerAttemptUserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_answer_attempt_user_id"`
	UserAnswerAttemptTargetType      int       `gorm:"not null" json:"user_answer_attempt_target_type"` // 1=reading, 2=quiz, 3=evaluation, 4=exam
	UserAnswerAttemptTargetID        int       `gorm:"not null" json:"user_answer_attempt_target_id"`
	UserAnswerAttemptQuestionID      int       `gorm:"not null" json:"user_answer_attempt_question_id"`
	UserAnswerAttemptAnswer          string    `gorm:"type:varchar(1);not null" json:"user_answer_attempt_answer"`
	UserAnswerAttemptIsCorrect       bool      `gorm:"default:false" json:"user_answer_attempt_is_correct"`
	UserAnswerAttemptCreatedAt       time.Time `gorm:"autoCreateTime" json:"user_answer_attempt_created_at"`

	UserAnswerAttemptBatchID     string    `gorm:"type:varchar(100)" json:"user_answer_attempt_batch_id"`
	UserAnswerAttemptSubmittedAt time.Time `json:"user_answer_attempt_submitted_at"`
}

func (UserAnswerAttempt) TableName() string {
	return "user_answer_attempts"
}
