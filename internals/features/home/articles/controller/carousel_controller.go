package controller

import (
	"log"
	"masjidku_backend/internals/features/home/articles/dto"
	"masjidku_backend/internals/features/home/articles/model"
	helper "masjidku_backend/internals/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CarouselController struct {
	DB *gorm.DB
}

func NewCarouselController(db *gorm.DB) *CarouselController {
	return &CarouselController{
		DB: db,
	}
}

// ✅ GET: Ambil semua carousel aktif (untuk publik)
func (ctrl *CarouselController) GetAllActiveCarousels(c *fiber.Ctx) error {
	var carousels []model.CarouselModel
	err := ctrl.DB.Preload("Article").
		Where("carousel_is_active = ?", true).
		Order("CASE WHEN carousel_order IS NOT NULL THEN 0 ELSE 1 END, carousel_order ASC, carousel_created_at DESC").
		Limit(3).
		Find(&carousels).Error
	if err != nil {
		log.Println("[ERROR] Gagal ambil data carousel:", err)
		return fiber.NewError(http.StatusInternalServerError, "Gagal ambil data carousel")
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil carousel",
		"data":    dto.ConvertCarouselListToDTO(carousels),
	})
}

// ✅ GET: Admin - Ambil semua carousel
func (ctrl *CarouselController) GetAllCarouselsAdmin(c *fiber.Ctx) error {
	var carousels []model.CarouselModel
	err := ctrl.DB.Preload("Article").
		Order("carousel_order").
		Find(&carousels).Error
	if err != nil {
		log.Println("[ERROR] Gagal ambil semua carousel admin:", err)
		return fiber.NewError(http.StatusInternalServerError, "Gagal ambil data")
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil data",
		"data":    dto.ConvertCarouselListToDTO(carousels),
	})
}

// ✅ POST: Admin - Tambah carousel
func (ctrl *CarouselController) CreateCarousel(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk tambah carousel")

	var req model.CarouselModel
	req.CarouselID = uuid.New()
	req.CarouselCreatedAt = time.Now()
	req.CarouselUpdatedAt = time.Now()

	// 📥 Ambil field selain gambar dari form
	req.CarouselTitle = c.FormValue("carousel_title")
	req.CarouselCaption = c.FormValue("carousel_caption")
	req.CarouselTargetURL = c.FormValue("carousel_target_url")
	req.CarouselType = c.FormValue("carousel_type")
	req.CarouselOrder, _ = strconv.Atoi(c.FormValue("carousel_order"))
	req.CarouselIsActive = c.FormValue("carousel_is_active") == "true"

	// ✨ Optional: carousel_article_id (UUID)
	if val := c.FormValue("carousel_article_id"); val != "" {
		articleID, err := uuid.Parse(val)
		if err == nil {
			req.CarouselArticleID = &articleID
		}
	}

	// 🖼️ Upload image
	fileHeader, err := c.FormFile("carousel_image_url")
	if err == nil && fileHeader != nil {
		// ✅ Jika file dikirim, upload
		url, err := helper.UploadImageAsWebPToSupabase("carousel", fileHeader)
		if err != nil {
			log.Println("[ERROR] Gagal upload gambar:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		req.CarouselImageURL = url
	} else {
		// ✅ Jika bukan file, ambil string
		url := c.FormValue("carousel_image_url")
		if url == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Gambar carousel wajib diisi (file atau URL)")
		}
		req.CarouselImageURL = url
	}

	// 💾 Simpan ke DB
	if err := ctrl.DB.Create(&req).Error; err != nil {
		log.Println("[ERROR] Gagal simpan carousel:", err)
		return fiber.NewError(http.StatusInternalServerError, "Gagal menyimpan data")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Carousel berhasil ditambahkan",
		"data":    dto.ConvertCarouselToDTO(req),
	})
}

func (ctrl *CarouselController) UpdateCarousel(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Update carousel ID: %s\n", id)

	var existing model.CarouselModel
	if err := ctrl.DB.Where("carousel_id = ?", id).First(&existing).Error; err != nil {
		log.Println("[ERROR] Carousel tidak ditemukan:", err)
		return fiber.NewError(fiber.StatusNotFound, "Data tidak ditemukan")
	}

	// Form parsing manual karena file + string campur
	carouselTitle := c.FormValue("carousel_title")
	carouselCaption := c.FormValue("carousel_caption")
	carouselTargetURL := c.FormValue("carousel_target_url")
	carouselType := c.FormValue("carousel_type")
	carouselOrder, _ := strconv.Atoi(c.FormValue("carousel_order"))
	carouselIsActive := c.FormValue("carousel_is_active") == "true"

	// Optional UUID untuk article
	var carouselArticleID *uuid.UUID
	if val := c.FormValue("carousel_article_id"); val != "" {
		if parsed, err := uuid.Parse(val); err == nil {
			carouselArticleID = &parsed
		}
	}

	// ⛔ Simpan URL gambar lama untuk hapus jika diganti
	oldImageURL := existing.CarouselImageURL
	newImageURL := oldImageURL // default: tidak berubah

	// 🔁 Cek apakah ada file baru
	fileHeader, err := c.FormFile("carousel_image_url")
	if err == nil && fileHeader != nil {
		// ✅ Jika file baru dikirim → upload
		url, err := helper.UploadImageAsWebPToSupabase("carousel", fileHeader)
		if err != nil {
			log.Println("[ERROR] Upload gambar gagal:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Upload gambar gagal")
		}
		newImageURL = url

		// 🔥 Hapus gambar lama jika ada
		if oldImageURL != "" && oldImageURL != newImageURL {
			bucket, path, err := helper.ExtractSupabasePath(oldImageURL)
			if err == nil {
				_ = helper.DeleteFromSupabase(bucket, path)
			}
		}
	} else {
		// 🔁 Coba ambil string URL jika ada
		str := c.FormValue("carousel_image_url")
		if str != "" {
			newImageURL = str
		}
	}

	// 📝 Update field
	existing.CarouselTitle = carouselTitle
	existing.CarouselCaption = carouselCaption
	existing.CarouselTargetURL = carouselTargetURL
	existing.CarouselType = carouselType
	existing.CarouselOrder = carouselOrder
	existing.CarouselIsActive = carouselIsActive
	existing.CarouselArticleID = carouselArticleID
	existing.CarouselImageURL = newImageURL
	existing.CarouselUpdatedAt = time.Now()

	// 💾 Simpan ke DB
	if err := ctrl.DB.Save(&existing).Error; err != nil {
		log.Println("[ERROR] Gagal update carousel:", err)
		return fiber.NewError(http.StatusInternalServerError, "Gagal update data")
	}

	return c.JSON(fiber.Map{
		"message": "Carousel berhasil diupdate",
		"data":    dto.ConvertCarouselToDTO(existing),
	})
}

// ✅ DELETE: Admin - Hapus carousel
func (ctrl *CarouselController) DeleteCarousel(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Where("carousel_id = ?", id).Delete(&model.CarouselModel{}).Error; err != nil {
		log.Println("[ERROR] Gagal hapus carousel:", err)
		return fiber.NewError(http.StatusInternalServerError, "Gagal hapus data")
	}

	return c.JSON(fiber.Map{
		"message": "Carousel berhasil dihapus",
	})
}
