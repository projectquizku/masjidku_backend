package dto

type TestExamResponse struct {
	TestExamID     uint   `json:"test_exam_id"`
	TestExamName   string `json:"test_exam_name"`
	TestExamStatus string `json:"test_exam_status"`
}
