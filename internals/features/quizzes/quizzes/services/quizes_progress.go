package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	userUnitModel "masjidku_backend/internals/features/lessons/units/model"
	quizzesModel "masjidku_backend/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SectionProgress struct {
	ID      uint `json:"id"`
	Score   int  `json:"score"`
	Attempt int  `json:"attempt"`
}
type UserQuizProgress struct {
	QuizID        uint `json:"quiz_id"`
	QuizAttempt   int  `json:"quiz_attempt"`
	QuizBestScore int  `json:"quiz_best_score"`
}

func UpdateUserSectionIfQuizCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	sectionQuizzesID uint,
	userQuizID uint,
	userQuizAttempt int,
	userQuizPercentageGrade int,
) error {
	log.Printf("[SERVICE] UpdateUserSectionIfQuizCompleted - user: %s, section: %d, quiz: %d", userID, sectionQuizzesID, userQuizID)

	// Ambil semua quiz dalam section
	var quizzesInSection []quizzesModel.QuizModel
	if err := db.Where("quiz_section_quizzes_id = ? AND deleted_at IS NULL", sectionQuizzesID).
		Find(&quizzesInSection).Error; err != nil {
		log.Println("[ERROR] Gagal mengambil daftar kuis dalam section:", err)
		return err
	}
	if len(quizzesInSection) == 0 {
		log.Println("[INFO] Section tidak memiliki quiz.")
		return nil
	}

	var allQuizIDsInSection = make(map[uint]struct{})
	for _, quiz := range quizzesInSection {
		allQuizIDsInSection[uint(quiz.QuizID)] = struct{}{}
	}

	// Ambil user_section_quizzes jika ada
	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?", userID, sectionQuizzesID).
		First(&userSection).Error

	newProgress := UserQuizProgress{
		QuizID:        userQuizID,
		QuizAttempt:   userQuizAttempt,
		QuizBestScore: userQuizPercentageGrade,
	}
	var progressList []UserQuizProgress

	if errors.Is(err, gorm.ErrRecordNotFound) {
		progressList = append(progressList, newProgress)
	} else if err == nil {
		if len(userSection.UserSectionQuizzesCompleteQuiz) > 0 {
			if err := json.Unmarshal(userSection.UserSectionQuizzesCompleteQuiz, &progressList); err != nil {
				log.Println("[ERROR] Gagal decode progress section sebelumnya:", err)
			}
		}
		found := false
		for i, p := range progressList {
			if p.QuizID == userQuizID {
				if userQuizAttempt > p.QuizAttempt {
					progressList[i].QuizAttempt = userQuizAttempt
				}
				if userQuizPercentageGrade > p.QuizBestScore {
					progressList[i].QuizBestScore = userQuizPercentageGrade
				}
				found = true
				break
			}
		}
		if !found {
			progressList = append(progressList, newProgress)
		}
	} else {
		log.Println("[ERROR] Gagal mengambil user_section_quizzes:", err)
		return err
	}

	// Cek kelengkapan kuis dan hitung grade
	completedQuizIDs := make(map[uint]bool)
	totalScore := 0
	for _, p := range progressList {
		completedQuizIDs[p.QuizID] = true
		totalScore += p.QuizBestScore
	}

	isAllQuizCompleted := len(completedQuizIDs) == len(allQuizIDsInSection)

	if isAllQuizCompleted {
		userSection.UserSectionQuizzesGradeResult = totalScore / len(progressList)
		log.Printf("[SERVICE] ✅ Semua kuis section lengkap. GradeResult = %d", userSection.UserSectionQuizzesGradeResult)
	} else {
		userSection.UserSectionQuizzesGradeResult = 0
		log.Println("[SERVICE] Kuis belum lengkap. GradeResult direset ke 0")
	}

	// Simpan atau buat baru user_section_quizzes
	if newJSON, err := json.Marshal(progressList); err == nil {
		userSection.UserSectionQuizzesUserID = userID
		userSection.UserSectionQuizzesSectionQuizzesID = sectionQuizzesID
		userSection.UserSectionQuizzesCompleteQuiz = datatypes.JSON(newJSON)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("[SERVICE] Membuat UserSectionQuizzes baru")
			return db.Create(&userSection).Error
		}
		log.Println("[SERVICE] Menyimpan update UserSectionQuizzes")
		return db.Save(&userSection).Error
	} else {
		log.Println("[ERROR] Gagal encode progress JSON:", err)
		return err
	}
}

func UpdateUserUnitIfSectionCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	unitID uint,
	completedSectionID uint,
) error {
	log.Printf("[SERVICE] UpdateUserUnitIfSectionCompleted - userID: %s, unitID: %d, completedSectionID: %d",
		userID.String(), unitID, completedSectionID)

	// 1. Cek progres user_section_quizzes
	var userSection quizzesModel.UserSectionQuizzesModel
	if err := db.Where(
		"user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?",
		userID, completedSectionID,
	).First(&userSection).Error; err != nil {
		log.Printf("[INFO] Section %d belum ada progress oleh user", completedSectionID)
		return nil
	}

	// 2. Ambil daftar quiz dalam section
	var section quizzesModel.SectionQuizzesModel
	if err := db.Preload("Quizzes").
		Where("section_quizzes_id = ?", completedSectionID).
		First(&section).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil section ID %d: %v", completedSectionID, err)
		return err
	}

	quizIDSet := make(map[uint]struct{})
	for _, quiz := range section.Quizzes {
		quizIDSet[uint(quiz.QuizID)] = struct{}{}
	}

	// 3. Decode progress kuis yang sudah dikerjakan
	var completedQuizData []struct {
		QuizID        int `json:"quiz_id"`
		QuizAttempt   int `json:"quiz_attempt"`
		QuizBestScore int `json:"quiz_best_score"`
	}
	if err := json.Unmarshal(userSection.UserSectionQuizzesCompleteQuiz, &completedQuizData); err != nil {
		log.Printf("[ERROR] Gagal decode complete_quiz: %v", err)
		return err
	}

	completedQuizIDs := make(map[uint]bool)
	for _, quiz := range completedQuizData {
		completedQuizIDs[uint(quiz.QuizID)] = true
	}

	// 4. Cek kelengkapan quiz
	for quizID := range quizIDSet {
		if !completedQuizIDs[quizID] {
			log.Printf("[INFO] Section %d belum lengkap, quiz ID %d belum dikerjakan", completedSectionID, quizID)
			return nil
		}
	}

	// 5. Ambil user_unit
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_unit_user_id = ? AND user_unit_unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil user_unit: %v", err)
		return err
	}

	// 6. Tambahkan completedSectionID jika belum ada
	var completedSectionIDs []int64
	if len(userUnit.UserUnitCompleteSectionQuizzes) > 0 {
		_ = json.Unmarshal(userUnit.UserUnitCompleteSectionQuizzes, &completedSectionIDs)
	}

	alreadyIncluded := false
	for _, sid := range completedSectionIDs {
		if uint(sid) == completedSectionID {
			alreadyIncluded = true
			break
		}
	}

	if !alreadyIncluded {
		completedSectionIDs = append(completedSectionIDs, int64(completedSectionID))
		if encoded, err := json.Marshal(completedSectionIDs); err == nil {
			userUnit.UserUnitCompleteSectionQuizzes = encoded
			userUnit.UpdatedAt = time.Now()
		}
	}

	// 7. Cek apakah semua section selesai
	var unit userUnitModel.UnitModel
	if err := db.Where("unit_id = ?", unitID).First(&unit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil unit: %v", err)
		return err
	}

	if len(unit.UnitTotalSectionQuizzes) > 0 && len(completedSectionIDs) == len(unit.UnitTotalSectionQuizzes) {
		totalQuizScore := 0
		sectionCount := 0

		for _, sectionID := range completedSectionIDs {
			var sec quizzesModel.UserSectionQuizzesModel
			if err := db.Where(
				"user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?",
				userID, sectionID,
			).First(&sec).Error; err != nil {
				log.Printf("[WARNING] Gagal ambil section progress: %v", err)
				continue
			}
			totalQuizScore += sec.UserSectionQuizzesGradeResult
			sectionCount++
		}

		if sectionCount > 0 {
			userUnit.UserUnitGradeQuiz = totalQuizScore / sectionCount
			userUnit.UserUnitGradeResult = (userUnit.UserUnitGradeQuiz + userUnit.UserUnitGradeExam + getGradeEvaluation(userUnit)) / 3
			userUnit.UserUnitIsPassed = userUnit.UserUnitGradeResult >= 70

			log.Printf("[SERVICE] ✅ Semua section selesai. GradeQuiz: %d, GradeResult: %d, IsPassed: %v",
				userUnit.UserUnitGradeQuiz, userUnit.UserUnitGradeResult, userUnit.UserUnitIsPassed)
		}
	}

	return db.Save(&userUnit).Error
}

func getGradeEvaluation(userUnit userUnitModel.UserUnitModel) int {
	type EvaluationAttempt struct {
		EvaluationAttemptCount int `json:"attempt"`
		EvaluationScore        int `json:"grade_evaluation"`
	}

	var evalData EvaluationAttempt
	if err := json.Unmarshal(userUnit.UserUnitAttemptEvaluation, &evalData); err != nil {
		log.Printf("[ERROR] Gagal mengurai JSON AttemptEvaluation: %v", err)
		return 0
	}
	return evalData.EvaluationScore
}
