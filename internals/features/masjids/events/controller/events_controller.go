package controller

import (
	"log"
	"masjidku_backend/internals/features/masjids/events/dto"
	"masjidku_backend/internals/features/masjids/events/model"
	masjidModel "masjidku_backend/internals/features/masjids/masjids/model"
	"strconv"
	"strings"
	"time"

	helper "masjidku_backend/internals/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	log.Println("[INFO] Menerima request untuk membuat event")

	// 🔒 Validasi user login (opsional)
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User belum login")
	}

	// 🔸 Ambil field wajib
	title := c.FormValue("event_title")
	description := c.FormValue("event_description")
	location := c.FormValue("event_location")
	masjidIDStr := c.FormValue("event_masjid_id")

	if title == "" || description == "" || masjidIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field wajib tidak lengkap")
	}

	// 🔸 Parse UUID masjid
	masjidID, err := uuid.Parse(masjidIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID masjid tidak valid")
	}

	// 🔸 Fallback lokasi ke nama masjid
	if location == "" {
		var masjid masjidModel.MasjidModel
		if err := ctrl.DB.First(&masjid, "masjid_id = ?", masjidID).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Gagal mengambil nama masjid untuk lokasi")
		}
		location = masjid.MasjidName
	}

	// 🔸 Upload gambar jika ada
	var imageURL *string
	if file, err := c.FormFile("event_image_url"); err == nil && file != nil {
		url, err := helper.UploadImageAsWebPToSupabase("event", file)
		if err != nil {
			log.Printf("[ERROR] Upload gambar gagal: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		imageURL = &url
	} else if val := c.FormValue("event_image_url"); val != "" {
		imageURL = &val
	}

	// 🔸 Parse opsional & tambahan metadata
	var (
		capacity               *int
		isPublic               = true
		isRegistrationRequired = false
		isPaid                 = false
		price                  *int
		paymentDeadline        *time.Time
	)

	if val := c.FormValue("event_capacity"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			capacity = &parsed
		}
	}
	if val := c.FormValue("event_is_public"); val == "false" {
		isPublic = false
	}
	if val := c.FormValue("event_is_registration_required"); val == "true" {
		isRegistrationRequired = true
	}
	if val := c.FormValue("event_is_paid"); val == "true" {
		isPaid = true
		if p := c.FormValue("event_price"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				price = &parsed
			}
		}
		if deadlineStr := c.FormValue("event_payment_deadline"); deadlineStr != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", deadlineStr); err == nil {
				paymentDeadline = &parsedTime
			}
		}
	}

	// 🔸 Susun DTO
	req := dto.EventRequest{
		EventTitle:                  title,
		EventDescription:            description,
		EventLocation:               location,
		EventImageURL:               imageURL,
		EventCapacity:               capacity,
		EventIsPublic:               isPublic,
		EventIsRegistrationRequired: isRegistrationRequired,
		EventIsPaid:                 isPaid,
		EventPrice:                  price,
		EventPaymentDeadline:        paymentDeadline,
		EventMasjidID:               masjidID,
	}

	// 🔸 Simpan ke DB
	newEvent := req.ToModel()
	if err := ctrl.DB.Create(newEvent).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan event: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan event")
	}

	// 🔸 Response sukses
	return c.Status(fiber.StatusCreated).JSON(dto.ToEventResponse(newEvent))
}

// 🟢 POST /api/a/events/by-masjid
func (ctrl *EventController) GetEventsByMasjid(c *fiber.Ctx) error {
	type Request struct {
		MasjidID string `json:"masjid_id"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil || body.MasjidID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Masjid ID tidak valid"})
	}

	var events []model.EventModel
	if err := ctrl.DB.
		Where("event_masjid_id = ?", body.MasjidID).
		Order("event_created_at DESC").
		Find(&events).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil data event: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil event",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil diambil",
		"data":    dto.ToEventResponseList(events),
	})
}

// 🟢 GET /api/a/events/all atau /api/u/events/all
func (ctrl *EventController) GetAllEvents(c *fiber.Ctx) error {
	var events []model.EventModel

	if err := ctrl.DB.Order("event_created_at DESC").Find(&events).Error; err != nil {
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

// 🟢 GET /api/u/events/:slug
func (ctrl *EventController) GetEventBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Slug tidak boleh kosong",
		})
	}

	var event model.EventModel
	if err := ctrl.DB.Where("event_slug = ?", slug).First(&event).Error; err != nil {
		log.Printf("[ERROR] Event dengan slug '%s' tidak ditemukan: %v", slug, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Event tidak ditemukan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil ditemukan",
		"data":    dto.ToEventResponse(&event),
	})
}

// 🟡 PUT /api/a/events/:id
func (ctrl *EventController) UpdateEvent(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk update event")

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID event tidak valid")
	}

	var existing model.EventModel
	if err := ctrl.DB.First(&existing, "event_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Event tidak ditemukan")
	}

	// 🔸 Field utama
	if title := c.FormValue("event_title"); title != "" {
		existing.EventTitle = title
		existing.EventSlug = strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	}
	if desc := c.FormValue("event_description"); desc != "" {
		existing.EventDescription = desc
	}
	if loc := c.FormValue("event_location"); loc != "" {
		existing.EventLocation = loc
	}
	if masjidIDStr := c.FormValue("event_masjid_id"); masjidIDStr != "" {
		if masjidID, err := uuid.Parse(masjidIDStr); err == nil {
			existing.EventMasjidID = masjidID
		}
	}

	// 🔁 Gambar baru
	if file, err := c.FormFile("event_image_url"); err == nil && file != nil {
		if existing.EventImageURL != nil {
			oldPath := helper.ExtractSupabaseStoragePath(*existing.EventImageURL)
			if oldPath != "" {
				_ = helper.DeleteFromSupabase("image", oldPath)
			}
		}
		url, err := helper.UploadImageAsWebPToSupabase("event", file)
		if err != nil {
			log.Println("[ERROR] Gagal upload gambar baru:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar baru")
		}
		existing.EventImageURL = &url
	} else if val := c.FormValue("event_image_url"); val != "" {
		if existing.EventImageURL != nil && !strings.HasPrefix(val, *existing.EventImageURL) {
			oldPath := helper.ExtractSupabaseStoragePath(*existing.EventImageURL)
			if oldPath != "" {
				_ = helper.DeleteFromSupabase("image", oldPath)
			}
		}
		existing.EventImageURL = &val
	}

	// 🔸 Field tambahan (opsional)
	if val := c.FormValue("event_capacity"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			existing.EventCapacity = &parsed
		}
	}

	if val := c.FormValue("event_is_public"); val != "" {
		existing.EventIsPublic = val != "false"
	}

	if val := c.FormValue("event_is_registration_required"); val != "" {
		existing.EventIsRegistrationRequired = val == "true"
	}

	if val := c.FormValue("event_is_paid"); val != "" {
		existing.EventIsPaid = val == "true"
	}

	if val := c.FormValue("event_price"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			existing.EventPrice = &parsed
		}
	}

	if val := c.FormValue("event_payment_deadline"); val != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
			existing.EventPaymentDeadline = &parsed
		}
	}

	// 🔄 Simpan
	if err := ctrl.DB.Save(&existing).Error; err != nil {
		log.Println("[ERROR] Gagal update event:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengupdate event")
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil diperbarui",
		"data":    dto.ToEventResponse(&existing),
	})
}

// 🔴 DELETE /api/a/events/:id
func (ctrl *EventController) DeleteEvent(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk menghapus event")

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID event tidak valid")
	}

	var event model.EventModel
	if err := ctrl.DB.First(&event, "event_id = ?", id).Error; err != nil {
		log.Printf("[ERROR] Event tidak ditemukan: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Event tidak ditemukan")
	}

	// ✅ Hapus gambar dari Supabase jika ada
	if event.EventImageURL != nil {
		oldPath := helper.ExtractSupabaseStoragePath(*event.EventImageURL)
		if oldPath != "" {
			_ = helper.DeleteFromSupabase("image", oldPath)
		}
	}

	// ✅ Hapus event dari DB
	if err := ctrl.DB.Delete(&event).Error; err != nil {
		log.Printf("[ERROR] Gagal menghapus event: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus event")
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil dihapus",
	})
}
