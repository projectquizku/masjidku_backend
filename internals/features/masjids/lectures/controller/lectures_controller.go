package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"masjidku_backend/internals/features/masjids/lectures/dto"
	"masjidku_backend/internals/features/masjids/lectures/model"
	helper "masjidku_backend/internals/helpers"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LectureController struct {
	DB *gorm.DB
}

func NewLectureController(db *gorm.DB) *LectureController {
	return &LectureController{DB: db}
}

func (ctrl *LectureController) CreateLecture(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk membuat kajian")

	// 🔍 Ambil field dari multipart form
	title := c.FormValue("lecture_title")
	description := c.FormValue("lecture_description")
	masjidIDStr := c.FormValue("lecture_masjid_id")
	totalSessionsStr := c.FormValue("total_lecture_sessions")
	statusStr := c.FormValue("lecture_status")
	certificateIDStr := c.FormValue("lecture_certificate_id")
	teachersJSON := c.FormValue("lecture_teachers")

	// Field tambahan
	isRegStr := c.FormValue("lecture_is_registration_required")
	isPaidStr := c.FormValue("lecture_is_paid")
	priceStr := c.FormValue("lecture_price")
	paymentDeadlineStr := c.FormValue("lecture_payment_deadline")
	paymentScope := c.FormValue("lecture_payment_scope")
	capacityStr := c.FormValue("lecture_capacity")
	isPublicStr := c.FormValue("lecture_is_public")

	// 📦 Parse UUID
	masjidID, err := uuid.Parse(masjidIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "lecture_masjid_id tidak valid"})
	}

	// 📦 Parse total session
	var totalSessions *int
	if totalSessionsStr != "" {
		tmp, err := strconv.Atoi(totalSessionsStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "total_lecture_sessions tidak valid"})
		}
		totalSessions = &tmp
	}

	// 📦 Parse status
	status := statusStr == "true"

	// 📦 Parse certificate_id
	var certificateID *uuid.UUID
	if certificateIDStr != "" {
		parsedID, err := uuid.Parse(certificateIDStr)
		if err == nil {
			certificateID = &parsedID
		}
	}

	// 📦 Parse teachers
	var teachers []dto.Teacher
	if teachersJSON != "" {
		if err := json.Unmarshal([]byte(teachersJSON), &teachers); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format lecture_teachers tidak valid"})
		}
	}

	// 🧱 Parse tambahan
	isReg := isRegStr == "true"
	isPaid := isPaidStr == "true"

	var price *int
	if priceStr != "" {
		p, err := strconv.Atoi(priceStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "lecture_price tidak valid"})
		}
		price = &p
	}

	var deadline *string
	if paymentDeadlineStr != "" {
		deadline = &paymentDeadlineStr
	}

	if paymentScope == "" {
		paymentScope = "lecture"
	}

	capacity := 0
	if capacityStr != "" {
		cap, err := strconv.Atoi(capacityStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "lecture_capacity tidak valid"})
		}
		capacity = cap
	}

	isPublic := isPublicStr != "false" // default true

	// 🖼️ Proses gambar
	var imageURL *string
	fileHeader, err := c.FormFile("lecture_image_url")
	if err == nil && fileHeader != nil {
		url, err := helper.UploadImageAsWebPToSupabase("lectures", fileHeader)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Gagal upload gambar: %v", err)})
		}
		imageURL = &url
	} else {
		urlStr := c.FormValue("lecture_image_url")
		if urlStr != "" {
			imageURL = &urlStr
		}
	}

	// 🧱 Bangun DTO
	req := dto.LectureRequest{
		LectureTitle:                  title,
		LectureDescription:            description,
		LectureMasjidID:               masjidID,
		TotalLectureSessions:          totalSessions,
		LectureImageURL:               imageURL,
		LectureTeachers:               teachers,
		LectureStatus:                 status,
		LectureCertificateID:          certificateID,
		LectureIsRegistrationRequired: isReg,
		LectureIsPaid:                 isPaid,
		LecturePrice:                  price,
		LecturePaymentDeadline:        deadline,
		LecturePaymentScope:           paymentScope,
		LectureCapacity:               capacity,
		LectureIsPublic:               isPublic,
	}

	// 💾 Simpan ke DB
	newLecture := req.ToModel()
	if err := ctrl.DB.Create(newLecture).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan data", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Kajian berhasil dibuat",
		"data":    dto.ToLectureResponse(newLecture),
	})
}

// ✅ POST /api/a/lectures/by-masjid-latest
func (ctrl *LectureController) GetByMasjidID(c *fiber.Ctx) error {
	type RequestBody struct {
		MasjidID string `json:"masjid_id"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Permintaan tidak valid",
		})
	}

	if body.MasjidID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Masjid ID wajib diisi",
		})
	}

	userIDRaw := c.Locals("user_id")

	// 🔍 Ambil semua lecture
	var lectures []model.LectureModel
	if err := ctrl.DB.
		Where("lecture_masjid_id = ?", body.MasjidID).
		Order("lecture_created_at DESC").
		Find(&lectures).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data lecture",
			"error":   err.Error(),
		})
	}

	if len(lectures) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Belum ada lecture untuk masjid ini",
		})
	}

	// 🧠 Persiapan response
	type Combined struct {
		Lecture     *dto.LectureResponse // <-- pakai pointer
		UserLecture *model.UserLectureModel
	}

	var response []Combined

	// 🧠 Siapkan lecture ID
	lectureIDs := make([]string, 0, len(lectures))
	for _, lec := range lectures {
		lectureIDs = append(lectureIDs, lec.LectureID.String())
	}

	// 🔄 Jika user login → ambil user_lecture
	userLectureMap := make(map[string]model.UserLectureModel)
	if userIDRaw != nil {
		userIDStr, ok := userIDRaw.(string)
		if ok && userIDStr != "" {
			userUUID, err := uuid.Parse(userIDStr)
			if err == nil {
				var userLectures []model.UserLectureModel
				if err := ctrl.DB.
					Where("user_lecture_user_id = ? AND user_lecture_lecture_id IN ?", userUUID, lectureIDs).
					Find(&userLectures).Error; err == nil {
					for _, ul := range userLectures {
						userLectureMap[ul.UserLectureLectureID.String()] = ul
					}
				}
			}
		}
	}

	// 🔀 Gabungkan lecture + user_lecture
	for _, lec := range lectures {
		item := Combined{
			Lecture:     dto.ToLectureResponse(&lec),
			UserLecture: nil,
		}

		if ul, ok := userLectureMap[lec.LectureID.String()]; ok {
			item.UserLecture = &ul
		}
		response = append(response, item)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil semua lecture masjid (dengan partisipasi user jika login)",
		"data":    response,
	})
}

// 🟢 GET /api/a/lectures/:id
func (ctrl *LectureController) GetLectureByID(c *fiber.Ctx) error {
	lectureIDStr := c.Params("id")
	lectureID, err := uuid.Parse(lectureIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID kajian tidak valid"})
	}

	var lecture model.LectureModel
	if err := ctrl.DB.First(&lecture, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Kajian tidak ditemukan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil detail kajian",
		"data":    dto.ToLectureResponse(&lecture),
	})
}

func (ctrl *LectureController) UpdateLecture(c *fiber.Ctx) error {
	lectureIDStr := c.Params("id")
	lectureID, err := uuid.Parse(lectureIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID kajian tidak valid"})
	}

	var lecture model.LectureModel
	if err := ctrl.DB.First(&lecture, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Kajian tidak ditemukan", "error": err.Error()})
	}

	// 🔍 Ambil field form
	title := c.FormValue("lecture_title")
	description := c.FormValue("lecture_description")
	totalSessionsStr := c.FormValue("total_lecture_sessions")
	statusStr := c.FormValue("lecture_status")
	certificateIDStr := c.FormValue("lecture_certificate_id")
	teachersJSON := c.FormValue("lecture_teachers")
	isPaidStr := c.FormValue("lecture_is_paid")
	isRegStr := c.FormValue("lecture_is_registration_required")
	priceStr := c.FormValue("lecture_price")
	scope := c.FormValue("lecture_payment_scope")
	deadlineStr := c.FormValue("lecture_payment_deadline")
	capacityStr := c.FormValue("lecture_capacity")
	isPublicStr := c.FormValue("lecture_is_public")

	// 🧠 Map update hanya jika field dikirim
	updated := make(map[string]interface{})

	if title != "" {
		updated["lecture_title"] = title
	}
	if description != "" {
		updated["lecture_description"] = description
	}
	if totalSessionsStr != "" {
		if val, err := strconv.Atoi(totalSessionsStr); err == nil {
			updated["total_lecture_sessions"] = val
		}
	}
	if statusStr != "" {
		updated["lecture_status"] = (statusStr == "true")
	}
	if certificateIDStr != "" {
		if id, err := uuid.Parse(certificateIDStr); err == nil {
			updated["lecture_certificate_id"] = id
		}
	}
	if teachersJSON != "" {
		var teachers []dto.Teacher
		if err := json.Unmarshal([]byte(teachersJSON), &teachers); err == nil {
			if data, err := json.Marshal(teachers); err == nil {
				updated["lecture_teachers"] = datatypes.JSON(data)
			}
		}
	}
	if isPaidStr != "" {
		updated["lecture_is_paid"] = (isPaidStr == "true")
	}
	if isRegStr != "" {
		updated["lecture_is_registration_required"] = (isRegStr == "true")
	}
	if isPublicStr != "" {
		updated["lecture_is_public"] = (isPublicStr != "false")
	}
	if priceStr != "" {
		if val, err := strconv.Atoi(priceStr); err == nil {
			updated["lecture_price"] = val
		}
	}
	if scope != "" {
		updated["lecture_payment_scope"] = scope
	}
	if deadlineStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", deadlineStr); err == nil {
			updated["lecture_payment_deadline"] = parsed
		}
	}
	if capacityStr != "" {
		if val, err := strconv.Atoi(capacityStr); err == nil {
			updated["lecture_capacity"] = val
		}
	}

	// 🖼️ Upload gambar baru jika ada
	if fileHeader, err := c.FormFile("lecture_image_url"); err == nil && fileHeader != nil {
		if url, err := helper.UploadImageAsWebPToSupabase("lectures", fileHeader); err == nil {
			updated["lecture_image_url"] = url
		} else {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Gagal upload gambar: %v", err)})
		}
	} else if urlStr := c.FormValue("lecture_image_url"); urlStr != "" {
		updated["lecture_image_url"] = urlStr
	}

	// 🧾 Jalankan update
	if err := ctrl.DB.Model(&lecture).Updates(updated).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memperbarui data", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Kajian berhasil diperbarui",
		"data":    dto.ToLectureResponse(&lecture),
	})
}

// 🔴 DELETE /api/a/lectures/:id
func (ctrl *LectureController) DeleteLecture(c *fiber.Ctx) error {
	lectureID := c.Params("id")

	if err := ctrl.DB.Delete(&model.LectureModel{}, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus kajian", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Kajian berhasil dihapus"})
}
