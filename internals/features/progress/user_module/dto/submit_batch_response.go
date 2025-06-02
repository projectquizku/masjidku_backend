package dto

import "time"

type SubmitUserResultRequest struct {
	Targets []SubmitTargetData `json:"targets"`
}

type SubmitTargetData struct {
	UserModuleAttemptTargetType      int                    `json:"user_module_attempt_target_type"`
	UserModuleAttemptTargetID        int                    `json:"user_module_attempt_target_id"`
	UserModuleAttemptPercentageGrade *int                   `json:"user_module_attempt_percentage_grade"`
	UserModuleAttemptTimeDuration    int                    `json:"user_module_attempt_time_duration"`
	UserModuleAttemptCreatedAt       time.Time              `json:"user_module_attempt_created_at"`
	Answers                          []AnswerAttemptRequest `json:"answers"`
}

type AnswerAttemptRequest struct {
	UserAnswerAttemptQuestionID      int       `json:"user_answer_attempt_question_id"`
	UserAnswerAttemptAnswer          string    `json:"user_answer_attempt_answer"`
	UserAnswerAttemptIsCorrect       bool      `json:"user_answer_attempt_is_correct"`
	UserAnswerAttemptCreatedAt       time.Time `json:"user_answer_attempt_created_at"`
}
