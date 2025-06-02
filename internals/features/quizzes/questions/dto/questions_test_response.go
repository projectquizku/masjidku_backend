package dto

type QuestionResponse struct {
	QuestionID            uint     `json:"question_id"`
	QuestionText          string   `json:"question_text"`
	QuestionAnswerChoices []string `json:"question_answer_choices"`
	QuestionCorrectAnswer string   `json:"question_correct_answer"`
	QuestionHelpParagraph string   `json:"question_paragraph_help"`
	QuestionExplanation   string   `json:"question_explanation"`
	QuestionAnswerText    string   `json:"question_answer_text"`
	QuestionStatus        string   `json:"question_status"`
}
