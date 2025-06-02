package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	themesOrLevelsModel "masjidku_backend/internals/features/lessons/themes_or_levels/model"
	userModel "masjidku_backend/internals/features/lessons/units/model"
	userSectionQuizzesModel "masjidku_backend/internals/features/quizzes/quizzes/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserUnitController struct {
	DB *gorm.DB
}

func NewUserUnitController(db *gorm.DB) *UserUnitController {
	return &UserUnitController{DB: db}
}

// 🟢 GET /api/user-units/:user_id
// Mengambil semua data progres unit milik user berdasarkan user_unit_user_id.
// Data yang dikembalikan termasuk relasi SectionProgress per unit.
func (ctrl *UserUnitController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	// Validasi UUID
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []userModel.UserUnitModel

	// Ambil semua user_unit berdasarkan user_unit_user_id
	if err := ctrl.DB.
		Preload("SectionProgress").
		Where("user_unit_user_id = ?", userID).
		Find(&data).Error; err != nil {

		log.Println("[ERROR] Gagal ambil data user_unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(data),
		"data":  data,
	})
}

// 🟢 GET /api/user-units/themes/:themes_or_levels_id
// Mengambil seluruh unit dalam sebuah theme (themes_or_levels_id) beserta progres user di tiap unit.
// Progress meliputi section_progress, complete_section_quizzes, dan field lain dari user_unit.
func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevels(c *fiber.Ctx) error {
	// 🔐 Ambil user_id dari token JWT
	userID, err := extractUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// 🔢 Ambil themes_or_levels_id dari URL
	themesID, err := strconv.Atoi(c.Params("themes_or_levels_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "themes_or_levels_id tidak valid"})
	}

	// 🔍 Cek user_theme
	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.
		Where("user_theme_user_id = ? AND user_theme_themes_or_level_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Data user_theme tidak ditemukan"})
	}

	// 📦 Ambil semua unit dan isinya
	var units []userModel.UnitModel
	if err := ctrl.DB.
		Preload("SectionQuizzes").
		Preload("SectionQuizzes.Quizzes").
		Where("unit_themes_or_level_id = ?", themesID).
		Find(&units).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil data unit"})
	}

	// 🔁 Mapping section → unit
	unitIDs, sectionQuizToUnit := extractUnitAndSectionMap(units)

	// 📈 Ambil progress per unit
	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.
		Where("user_unit_user_id = ? AND user_unit_unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil data progress unit"})
	}

	// 📈 Ambil progress per section
	var sectionProgressList []userSectionQuizzesModel.UserSectionQuizzesModel
	if err := ctrl.DB.
		Where("user_section_quizzes_user_id = ?", userID).
		Where("user_section_quizzes_section_quizzes_id IN ?", keys(sectionQuizToUnit)).
		Find(&sectionProgressList).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil section_progress user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil section_progress"})
	}

	// 🔁 Organisasi ke unit masing-masing
	progressPerUnit := make(map[uint][]userSectionQuizzesModel.UserSectionQuizzesModel)
	completedMap := make(map[uint]pq.Int64Array)
	for _, sp := range sectionProgressList {
		unitID := sectionQuizToUnit[sp.UserSectionQuizzesSectionQuizzesID]
		progressPerUnit[unitID] = append(progressPerUnit[unitID], sp)
		if len(sp.UserSectionQuizzesCompleteQuiz) > 0 {
			completedMap[unitID] = append(completedMap[unitID], int64(sp.UserSectionQuizzesSectionQuizzesID))
		}
	}

	// 💡 Gabungkan progress ke user_unit
	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		unitID := u.UserUnitUnitID
		u.SectionProgress = progressPerUnit[unitID]

		if completed, ok := completedMap[u.UserUnitUnitID]; ok && len(completed) > 0 {
			// 🔄 Marshal array ke JSON
			completedJSON, err := json.Marshal(completed)
			if err == nil {
				u.UserUnitCompleteSectionQuizzes = datatypes.JSON(completedJSON)

				// 💾 Update ke DB
				_ = ctrl.DB.Model(&userModel.UserUnitModel{}).
					Where("user_unit_id = ?", u.UserUnitID).
					Update("user_unit_complete_section_quizzes", datatypes.JSON(completedJSON)).Error
			}
		}

		progressMap[unitID] = u
	}

	// 📤 Bentuk response akhir
	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}
	var result []ResponseUnit
	for _, unit := range units {
		userProgress := progressMap[unit.UnitID]
		result = append(result, ResponseUnit{
			UnitModel:    unit,
			UserProgress: userProgress,
		})
	}

	return c.JSON(fiber.Map{"data": result})
}
func extractUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	val := c.Locals("user_id")
	if val == nil {
		return uuid.Nil, fmt.Errorf("Unauthorized - user_id tidak ditemukan dalam token")
	}
	str, ok := val.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("Unauthorized - format user_id tidak valid")
	}
	return uuid.Parse(str)
}

func keys(m map[uint]uint) []uint {
	k := make([]uint, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func extractUnitAndSectionMap(units []userModel.UnitModel) ([]uint, map[uint]uint) {
	var unitIDs []uint
	sectionMap := make(map[uint]uint)
	for _, unit := range units {
		unitIDs = append(unitIDs, unit.UnitID)
		for _, section := range unit.SectionQuizzes {
			sectionMap[section.SectionQuizzesID] = unit.UnitID
		}
	}
	return unitIDs, sectionMap
}
