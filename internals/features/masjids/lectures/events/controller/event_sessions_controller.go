package controller

import (
	"log"
	"masjidku_backend/internals/features/masjids/lectures/events/dto"
	"masjidku_backend/internals/features/masjids/lectures/events/model"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventSessionController struct {
	DB *gorm.DB
}

func NewEventSessionController(db *gorm.DB) *EventSessionController {
	return &EventSessionController{DB: db}
}

// 🟢 POST /api/a/event-sessions
// 🟢 POST /api/a/event-sessions
func (ctrl *EventSessionController) CreateEventSession(c *fiber.Ctx) error {
	// Ambil user_id dari token (middleware harus sudah set ini di Locals)
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		log.Println("[ERROR] Gagal mendapatkan user_id dari token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak terautentikasi",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[ERROR] Gagal parsing user_id: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User ID tidak valid",
		})
	}

	var req dto.EventSessionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Body parser gagal: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Permintaan tidak valid",
			"error":   err.Error(),
		})
	}

	session := req.ToModel()
	session.EventSessionCreatedBy = &userID // ✅ Set created_by dari token

	if err := ctrl.DB.Create(session).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan event session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan event session",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event session berhasil dibuat",
		"data":    dto.ToEventSessionResponse(session),
	})
}

// 🟢 GET /api/u/event-sessions/by-event/:event_id
func (ctrl *EventSessionController) GetEventSessionsByEvent(c *fiber.Ctx) error {
	eventID := c.Params("event_id")
	if eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event ID tidak boleh kosong",
		})
	}

	var sessions []model.EventSessionModel
	if err := ctrl.DB.Where("event_session_event_id = ?", eventID).
		Order("event_session_start_time ASC").
		Find(&sessions).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil event sessions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil event sessions",
			"error":   err.Error(),
		})
	}

	now := time.Now()
	var result []dto.EventSessionResponse
	for _, s := range sessions {
		status := "upcoming"
		if now.After(s.EventSessionStartTime) && now.Before(s.EventSessionEndTime) {
			status = "ongoing"
		} else if now.After(s.EventSessionEndTime) {
			status = "completed"
		}

		resp := dto.ToEventSessionResponse(&s)
		resp.EventSessionStatus = status
		result = append(result, *resp)
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil event sessions",
		"data":    result,
	})
}

// 🟢 GET /api/u/event-sessions/all
func (ctrl *EventSessionController) GetAllEventSessions(c *fiber.Ctx) error {
	var sessions []model.EventSessionModel
	if err := ctrl.DB.Order("event_session_start_time DESC").Find(&sessions).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil semua event session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data event sessions",
			"error":   err.Error(),
		})
	}

	now := time.Now()
	var result []dto.EventSessionResponse
	for _, s := range sessions {
		status := "upcoming"
		if now.After(s.EventSessionStartTime) && now.Before(s.EventSessionEndTime) {
			status = "ongoing"
		} else if now.After(s.EventSessionEndTime) {
			status = "completed"
		}

		resp := dto.ToEventSessionResponse(&s)
		resp.EventSessionStatus = status
		result = append(result, *resp)
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil semua event session",
		"data":    result,
	})
}

// 🟢 GET /api/u/event-sessions/upcoming
func (ctrl *EventSessionController) GetUpcomingEventSessions(c *fiber.Ctx) error {
	var sessions []model.EventSessionModel

	if err := ctrl.DB.
		Where("event_session_start_time > ? AND event_session_is_public = ?", time.Now(), true).
		Order("event_session_start_time ASC").
		Find(&sessions).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil sesi event upcoming: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil sesi event upcoming",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil sesi event yang akan datang",
		"data":    dto.ToEventSessionResponseList(sessions),
	})
}
