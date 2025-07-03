package controller

import (
	"log"
	"masjidku_backend/internals/features/masjids/events/dto"
	"masjidku_backend/internals/features/masjids/events/model"
	helper "masjidku_backend/internals/helpers"

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
func (ctrl *EventSessionController) CreateEventSession(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk membuat sesi event")

	// ✅ Ambil user_id dari token
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "User tidak terautentikasi")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID tidak valid")
	}

	// ✅ Ambil form values
	get := func(key string) string {
		return c.FormValue(key)
	}

	title := get("event_session_title")
	desc := get("event_session_description")
	startStr := get("event_session_start_time")
	endStr := get("event_session_end_time")
	masjidIDStr := get("event_session_masjid_id")
	eventIDStr := get("event_session_event_id")

	// ✅ Validasi wajib
	if title == "" || desc == "" || startStr == "" || endStr == "" || masjidIDStr == "" || eventIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field wajib tidak lengkap")
	}

	// ✅ Parsing waktu
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format waktu mulai tidak valid (RFC3339)")
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format waktu selesai tidak valid (RFC3339)")
	}
	if endTime.Before(startTime) {
		return fiber.NewError(fiber.StatusBadRequest, "Waktu selesai harus setelah waktu mulai")
	}

	// ✅ Parsing UUID
	masjidID, err := uuid.Parse(masjidIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Masjid ID tidak valid")
	}
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Event ID tidak valid")
	}

	// Opsional
	location := get("event_session_location")

	// ✅ Gambar
	var imageURL *string
	if file, err := c.FormFile("event_session_image_url"); err == nil && file != nil {
		url, err := helper.UploadImageAsWebPToSupabase("event_sessions", file)
		if err != nil {
			log.Println("[ERROR] Gagal upload gambar:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar sesi")
		}
		imageURL = &url
	} else if val := get("event_session_image_url"); val != "" {
		imageURL = &val
	}

	// ✅ Susun request DTO
	req := dto.EventSessionRequest{
		EventSessionEventID:     eventID,
		EventSessionTitle:       title,
		EventSessionDescription: desc,
		EventSessionStartTime:   startTime,
		EventSessionEndTime:     endTime,
		EventSessionLocation:    location,
		EventSessionImageURL:    "",
		EventSessionMasjidID:    masjidID,
		EventSessionCreatedBy:   &userID,
	}
	if imageURL != nil {
		req.EventSessionImageURL = *imageURL
	}

	// ✅ Simpan ke DB
	model := req.ToModel()
	if err := ctrl.DB.Create(model).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan sesi event: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan sesi event")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sesi event berhasil dibuat",
		"data":    dto.ToEventSessionResponse(model),
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

	// 🔸 Local struct extend DTO
	type ExtendedResponse struct {
		dto.EventSessionResponse
		EventSessionStatus string `json:"event_session_status"`
	}

	now := time.Now()
	var result []ExtendedResponse
	for _, s := range sessions {
		status := "upcoming"
		if now.After(s.EventSessionStartTime) && now.Before(s.EventSessionEndTime) {
			status = "ongoing"
		} else if now.After(s.EventSessionEndTime) {
			status = "completed"
		}

		resp := dto.ToEventSessionResponse(&s)
		result = append(result, ExtendedResponse{
			EventSessionResponse: *resp,
			EventSessionStatus:   status,
		})
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

	// 🔸 Extend response dengan status
	type ExtendedResponse struct {
		dto.EventSessionResponse
		EventSessionStatus string `json:"event_session_status"`
	}

	now := time.Now()
	var result []ExtendedResponse
	for _, s := range sessions {
		status := "upcoming"
		if now.After(s.EventSessionStartTime) && now.Before(s.EventSessionEndTime) {
			status = "ongoing"
		} else if now.After(s.EventSessionEndTime) {
			status = "completed"
		}

		resp := dto.ToEventSessionResponse(&s)
		result = append(result, ExtendedResponse{
			EventSessionResponse: *resp,
			EventSessionStatus:   status,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil semua event session",
		"data":    result,
	})
}

// 🟢 GET /api/u/event-sessions/upcoming/:masjid_id?
func (ctrl *EventSessionController) GetUpcomingEventSessions(c *fiber.Ctx) error {
	var sessions []model.EventSessionModel

	masjidIDStr := c.Params("masjid_id")
	query := ctrl.DB.
		Where("event_session_start_time > ?", time.Now())

	if masjidIDStr != "" {
		masjidID, err := uuid.Parse(masjidIDStr)
		if err != nil {
			log.Printf("[ERROR] Invalid masjid_id format from path: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Format ID masjid tidak valid",
				"error":   "Invalid UUID format for masjid_id in path",
			})
		}
		query = query.Where("event_session_masjid_id = ?", masjidID)
	} else {
		log.Printf("[INFO] GetUpcomingEventSessions dipanggil tanpa masjid_id di path. Mengambil semua sesi.")
	}

	if err := query.Order("event_session_start_time ASC").Find(&sessions).Error; err != nil {
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

// 🟡 PUT /api/a/event-sessions/:id
func (ctrl *EventSessionController) UpdateEventSession(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk update sesi event")

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID sesi event tidak valid")
	}

	var session model.EventSessionModel
	if err := ctrl.DB.First(&session, "event_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Sesi event tidak ditemukan")
	}

	get := func(key string) string { return c.FormValue(key) }

	// Update field sesuai yang dikirim
	if val := get("event_session_title"); val != "" {
		session.EventSessionTitle = val
	}
	if val := get("event_session_description"); val != "" {
		session.EventSessionDescription = val
	}
	if val := get("event_session_location"); val != "" {
		session.EventSessionLocation = val
	}

	// Time
	if val := get("event_session_start_time"); val != "" {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			session.EventSessionStartTime = t
		}
	}
	if val := get("event_session_end_time"); val != "" {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			session.EventSessionEndTime = t
		}
	}

	// Relasi
	if val := get("event_session_event_id"); val != "" {
		if uuidVal, err := uuid.Parse(val); err == nil {
			session.EventSessionEventID = uuidVal
		}
	}
	if val := get("event_session_masjid_id"); val != "" {
		if uuidVal, err := uuid.Parse(val); err == nil {
			session.EventSessionMasjidID = uuidVal
		}
	}

	// ✅ Gambar: jika diganti, hapus yang lama
	if file, err := c.FormFile("event_session_image_url"); err == nil && file != nil {
		if session.EventSessionImageURL != "" {
			oldPath := helper.ExtractSupabaseStoragePath(session.EventSessionImageURL)
			_ = helper.DeleteFromSupabase("image", oldPath)
		}
		newURL, err := helper.UploadImageAsWebPToSupabase("event_sessions", file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		session.EventSessionImageURL = newURL
	} else if val := get("event_session_image_url"); val != "" {
		if val != session.EventSessionImageURL {
			if session.EventSessionImageURL != "" {
				oldPath := helper.ExtractSupabaseStoragePath(session.EventSessionImageURL)
				_ = helper.DeleteFromSupabase("image", oldPath)
			}
			session.EventSessionImageURL = val
		}
	}

	// Simpan perubahan
	if err := ctrl.DB.Save(&session).Error; err != nil {
		log.Println("[ERROR] Gagal update sesi event:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal update sesi event")
	}

	return c.JSON(fiber.Map{
		"message": "Sesi event berhasil diperbarui",
		"data":    dto.ToEventSessionResponse(&session),
	})
}

// 🔴 DELETE /api/a/event-sessions/:id
func (ctrl *EventSessionController) DeleteEventSession(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk menghapus sesi event")

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID sesi event tidak valid")
	}

	var session model.EventSessionModel
	if err := ctrl.DB.First(&session, "event_session_id = ?", id).Error; err != nil {
		log.Printf("[ERROR] Sesi event tidak ditemukan: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Sesi event tidak ditemukan")
	}

	// ✅ Hapus gambar dari Supabase jika ada
	if session.EventSessionImageURL != "" {
		oldPath := helper.ExtractSupabaseStoragePath(session.EventSessionImageURL)
		if oldPath != "" {
			_ = helper.DeleteFromSupabase("image", oldPath)
		}
	}

	// ✅ Hapus dari database
	if err := ctrl.DB.Delete(&session).Error; err != nil {
		log.Printf("[ERROR] Gagal menghapus sesi event: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus sesi event")
	}

	return c.JSON(fiber.Map{
		"message": "Sesi event berhasil dihapus",
	})
}
