package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	activityService "masjidku_backend/internals/features/progress/daily_activities/service"

	"masjidku_backend/internals/features/progress/user_module/dto"
	"masjidku_backend/internals/features/progress/user_module/model"
	evaluationService "masjidku_backend/internals/features/quizzes/evaluations/service"
	examService "masjidku_backend/internals/features/quizzes/exams/service"
	quizService "masjidku_backend/internals/features/quizzes/quizzes/services"
	readingService "masjidku_backend/internals/features/quizzes/readings/service"
)

type SubmitBatchController struct {
	DB *gorm.DB
}

func NewSubmitBatchController(db *gorm.DB) *SubmitBatchController {
	return &SubmitBatchController{DB: db}
}

func (sbc *SubmitBatchController) Submit(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user_id tidak valid",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user_id bukan UUID valid",
		})
	}

	submittedAt := time.Now().UTC()
	batchID := fmt.Sprintf("batch-%s-%s", userID.String(), submittedAt.Format("20060102T150405"))

	var payload dto.SubmitUserResultRequest
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("[ERROR] Body parse failed: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	moduleAttempts := make([]model.UserModuleAttempt, 0)
	answerAttempts := make([]model.UserAnswerAttempt, 0)

	// Iterate over the targets and prepare the attempts without UserModuleAttemptAttemptNumber
	for _, target := range payload.Targets {
		moduleAttempts = append(moduleAttempts, model.UserModuleAttempt{
			UserModuleAttemptUserID:          userID,
			UserModuleAttemptTargetType:      target.UserModuleAttemptTargetType,
			UserModuleAttemptTargetID:        target.UserModuleAttemptTargetID,
			UserModuleAttemptPercentageGrade: target.UserModuleAttemptPercentageGrade,
			UserModuleAttemptTimeDuration:    &target.UserModuleAttemptTimeDuration,
			UserModuleAttemptCreatedAt:       target.UserModuleAttemptCreatedAt,
			UserModuleAttemptBatchID:         batchID,
			UserModuleAttemptSubmittedAt:     submittedAt,
		})

		// Handle the answer attempts as well
		for _, ans := range target.Answers {
			answerAttempts = append(answerAttempts, model.UserAnswerAttempt{
				UserAnswerAttemptUserID:      userID,
				UserAnswerAttemptTargetType:  target.UserModuleAttemptTargetType,
				UserAnswerAttemptTargetID:    target.UserModuleAttemptTargetID,
				UserAnswerAttemptQuestionID:  ans.UserAnswerAttemptQuestionID,
				UserAnswerAttemptAnswer:      ans.UserAnswerAttemptAnswer,
				UserAnswerAttemptIsCorrect:   ans.UserAnswerAttemptIsCorrect,
				UserAnswerAttemptCreatedAt:   ans.UserAnswerAttemptCreatedAt,
				UserAnswerAttemptBatchID:     batchID,
				UserAnswerAttemptSubmittedAt: submittedAt,
			})
		}
	}

	// Save attempts to DB in a transaction
	err = sbc.DB.Transaction(func(tx *gorm.DB) error {
		if len(moduleAttempts) > 0 {
			if err := tx.Create(&moduleAttempts).Error; err != nil {
				return err
			}
		}
		if len(answerAttempts) > 0 {
			if err := tx.Create(&answerAttempts).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Failed to save attempts for user_id=%s: %v", userID.String(), err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}

	// Handle points and progress
	for _, target := range payload.Targets {
		// Calculate attempts dynamically based on progress, no need to track it manually
		attempt := 1 // Default to 1 if there's no need for attempt tracking

		switch target.UserModuleAttemptTargetType {
		case 1:
			_ = readingService.UpdateUserUnitFromReading(sbc.DB, userID, uint(target.UserModuleAttemptTargetID))
			_ = readingService.AddPointFromReading(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), attempt)
			_ = activityService.UpdateOrInsertDailyActivity(sbc.DB, userID)
		case 2:
			_ = quizService.AddPointFromQuiz(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), attempt)
			_ = activityService.UpdateOrInsertDailyActivity(sbc.DB, userID)
			_ = quizService.UpdateUserSectionIfQuizCompleted(sbc.DB, userID, 0, uint(target.UserModuleAttemptTargetID), attempt, GradeOrZero(target.UserModuleAttemptPercentageGrade))
			_ = quizService.UpdateUserUnitIfSectionCompleted(sbc.DB, userID, 0, 0)
		case 3:
			grade := GradeOrZero(target.UserModuleAttemptPercentageGrade)
			_ = evaluationService.UpdateUserUnitFromEvaluation(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), grade)
			_ = evaluationService.AddPointFromEvaluation(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), attempt)
			_ = activityService.UpdateOrInsertDailyActivity(sbc.DB, userID)
		case 4:
			_ = examService.AddPointFromExam(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), attempt)
			_ = activityService.UpdateOrInsertDailyActivity(sbc.DB, userID)
			_ = examService.UpdateUserUnitFromExam(sbc.DB, userID, uint(target.UserModuleAttemptTargetID), GradeOrZero(target.UserModuleAttemptPercentageGrade))
		}
	}

	log.Printf("[SUCCESS] Submit completed for user_id=%s batch_id=%s", userID.String(), batchID)
	return c.JSON(fiber.Map{
		"status":       true,
		"message":      "Submit success",
		"batch_id":     batchID,
		"submitted_at": submittedAt,
	})
}

func GradeOrZero(grade *int) int {
	if grade != nil {
		return *grade
	}
	return 0
}

// 👨‍💼 Endpoint admin melihat semua module attempts
func (sbc *SubmitBatchController) GetAllModuleAttempts(c *fiber.Ctx) error {
	var attempts []model.UserModuleAttempt
	if err := sbc.DB.Order("user_module_attempt_created_at desc").Find(&attempts).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch module attempts: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve module attempts"})
	}
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "All user module attempts retrieved",
		"data":    attempts,
	})
}

func (sbc *SubmitBatchController) GetBatch(c *fiber.Ctx) error {
	log.Println("[DEBUG] Handler GET /get-batch/:batch_id DIPANGGIL")
	batchID := c.Params("batch_id")

	var attempts []model.UserModuleAttempt
	if err := sbc.DB.Where("user_module_attempt_batch_id = ?", batchID).Find(&attempts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch batch data"})
	}

	var answers []model.UserAnswerAttempt
	if err := sbc.DB.Where("user_answer_attempt_batch_id = ?", batchID).Find(&answers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch answer data"})
	}

	return c.JSON(fiber.Map{
		"batch_id":     batchID,
		"submitted_at": attempts[0].UserModuleAttemptSubmittedAt, // asumsi batch sama waktu
		"attempts":     attempts,
		"answers":      answers,
	})
}

// 👨‍💼 Endpoint admin melihat semua answer attempts
func (sbc *SubmitBatchController) GetAllAnswerAttempts(c *fiber.Ctx) error {
	var answers []model.UserAnswerAttempt
	if err := sbc.DB.Order("user_answer_attempt_created_at desc").Find(&answers).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch answer attempts: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve answer attempts"})
	}
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "All user answer attempts retrieved",
		"data":    answers,
	})
}
