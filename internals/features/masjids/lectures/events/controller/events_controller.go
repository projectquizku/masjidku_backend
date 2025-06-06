package controller

import (
	"log"
	"masjidku_backend/internals/features/masjids/lectures/events/dto"
	"masjidku_backend/internals/features/masjids/lectures/events/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EventController struct {
	DB *gorm.DB
}

func NewEventController(db *gorm.DB) *EventController {
	return &EventController{DB: db}
}

// 🟢 POST /api/a/events
func (ctrl *EventController) CreateEvent(c *fiber.Ctx) error {
	var req dto.EventRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Body parser gagal: %v", err)
		return c.Status(400).JSON(fiber.Map{"message": "Permintaan tidak valid", "error": err.Error()})
	}

	newEvent := req.ToModel()
	if err := ctrl.DB.Create(newEvent).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan event: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan event", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil ditambahkan",
		"data":    dto.ToEventResponse(newEvent),
	})
}

// 🟢 POST /api/a/events/by-masjid
func (ctrl *EventController) GetEventsByMasjid(c *fiber.Ctx) error {
	type Request struct {
		MasjidID string `json:"masjid_id"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil || body.MasjidID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Masjid ID tidak valid"})
	}

	var events []model.EventModel
	if err := ctrl.DB.
		Where("event_masjid_id = ?", body.MasjidID).
		Order("event_start_time ASC").
		Find(&events).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil data event: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil event", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil diambil",
		"data":    dto.ToEventResponseList(events),
	})
}

// 🟢 GET /api/a/events/all atau /api/u/events/all
func (ctrl *EventController) GetAllEvents(c *fiber.Ctx) error {
	var events []model.EventModel

	if err := ctrl.DB.Order("event_start_time desc").Find(&events).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil semua event: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data event",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil semua event",
		"data":    dto.ToEventResponseList(events),
	})
}
