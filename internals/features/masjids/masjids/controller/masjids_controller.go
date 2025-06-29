package controller

import (
	"fmt"
	"log"
	"net/http"

	"masjidku_backend/internals/features/masjids/masjids/dto"
	"masjidku_backend/internals/features/masjids/masjids/model"
	helper "masjidku_backend/internals/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasjidController struct {
	DB *gorm.DB
}

func NewMasjidController(db *gorm.DB) *MasjidController {
	return &MasjidController{DB: db}
}

// 🟢 GET ALL MASJIDS
func (mc *MasjidController) GetAllMasjids(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all masjids")

	var masjids []model.MasjidModel
	if err := mc.DB.Find(&masjids).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch masjids: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data masjid",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d masjids\n", len(masjids))

	// 🔁 Transform ke DTO
	var masjidDTOs []dto.MasjidResponse
	for _, m := range masjids {
		masjidDTOs = append(masjidDTOs, dto.FromModelMasjid(&m))
	}

	return c.JSON(fiber.Map{
		"message": "Data semua masjid berhasil diambil",
		"total":   len(masjidDTOs),
		"data":    masjidDTOs,
	})
}

// 🟢 GET VERIFIED MASJIDS
func (mc *MasjidController) GetAllVerifiedMasjids(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all verified masjids")

	var masjids []model.MasjidModel
	if err := mc.DB.Where("masjid_is_verified = ?", true).Find(&masjids).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch verified masjids: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data masjid terverifikasi",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d verified masjids\n", len(masjids))

	// 🔁 Transform ke DTO
	var masjidDTOs []dto.MasjidResponse
	for _, m := range masjids {
		masjidDTOs = append(masjidDTOs, dto.FromModelMasjid(&m))
	}

	return c.JSON(fiber.Map{
		"message": "Data masjid terverifikasi berhasil diambil",
		"total":   len(masjidDTOs),
		"data":    masjidDTOs,
	})
}

// 🟢 GET MASJID BY SLUG
func (mc *MasjidController) GetMasjidBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	log.Printf("[INFO] Fetching masjid with slug: %s\n", slug)

	var masjid model.MasjidModel
	if err := mc.DB.Where("masjid_slug = ?", slug).First(&masjid).Error; err != nil {
		log.Printf("[ERROR] Masjid with slug %s not found\n", slug)
		return c.Status(404).JSON(fiber.Map{
			"error": "Masjid tidak ditemukan",
		})
	}

	log.Printf("[SUCCESS] Retrieved masjid: %s\n", masjid.MasjidName)

	// 🔁 Transform ke DTO
	masjidDTO := dto.FromModelMasjid(&masjid)

	return c.JSON(fiber.Map{
		"message": "Data masjid berhasil diambil",
		"data":    masjidDTO,
	})
}

// 🔰 Create Masjid
func (mc *MasjidController) CreateMasjid(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk membuat masjid")

	// 📥 Ambil input dari form-data
	masjidName := c.FormValue("masjid_name")
	masjidBioShort := c.FormValue("masjid_bio_short")
	masjidLocation := c.FormValue("masjid_location")
	latitude := c.FormValue("masjid_latitude")
	longitude := c.FormValue("masjid_longitude")
	slug := helper.GenerateSlug(masjidName)
	isVerified := c.FormValue("masjid_is_verified") == "true"
	instagramURL := c.FormValue("masjid_instagram_url")
	whatsappURL := c.FormValue("masjid_whatsapp_url")
	youtubeURL := c.FormValue("masjid_youtube_url")

	// 🖼️ Upload gambar
	// 🖼️ Proses masjid_image_url: bisa file ATAU string
	var imageURL string

	// Coba ambil file
	fileHeader, err := c.FormFile("masjid_image_url")
	if err == nil && fileHeader != nil {
		// ✅ Jika file dikirim, upload ke Supabase
		url, err := helper.UploadImageAsWebPToSupabase("masjid", fileHeader)
		if err != nil {
			log.Printf("[ERROR] Upload gambar gagal: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		imageURL = url
	} else {
		// ✅ Coba ambil string URL
		strURL := c.FormValue("masjid_image_url")
		if strURL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Masjid image wajib dikirim sebagai file atau URL string"})
		}
		imageURL = strURL
	}

	// 🧭 Parse koordinat
	lat, lng, err := helper.ParseCoordinates(latitude, longitude)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Koordinat tidak valid"})
	}

	// 🧱 Siapkan model
	masjidID := uuid.New()
	newMasjid := &model.MasjidModel{
		MasjidID:           masjidID,
		MasjidName:         masjidName,
		MasjidBioShort:     masjidBioShort,
		MasjidLocation:     masjidLocation,
		MasjidLatitude:     lat,
		MasjidLongitude:    lng,
		MasjidImageURL:     imageURL,
		MasjidSlug:         slug,
		MasjidIsVerified:   isVerified,
		MasjidInstagramURL: instagramURL,
		MasjidWhatsappURL:  whatsappURL,
		MasjidYoutubeURL:   youtubeURL,
	}

	// 💾 Simpan ke DB
	if err := mc.DB.Create(newMasjid).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan masjid: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan data masjid"})
	}

	// 📤 Kirim response
	return c.Status(http.StatusCreated).JSON(dto.FromModelMasjid(newMasjid))
}

// 🟢 UPDATE MASJID (Partial Update)
func (mc *MasjidController) UpdateMasjid(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating masjid with ID: %s\n", id)

	masjidUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format ID tidak valid"})
	}

	var existing model.MasjidModel
	if err := mc.DB.First(&existing, "masjid_id = ?", masjidUUID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Masjid tidak ditemukan"})
	}

	// 🧭 Koordinat (opsional)
	latitude := c.FormValue("masjid_latitude")
	longitude := c.FormValue("masjid_longitude")
	if latitude != "" && longitude != "" {
		lat, lng, err := helper.ParseCoordinates(latitude, longitude)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Koordinat tidak valid"})
		}
		existing.MasjidLatitude = lat
		existing.MasjidLongitude = lng
	}

	// 🖼️ Gambar (string atau file)
	fileHeader, errFile := c.FormFile("masjid_image_url")
	if errFile == nil && fileHeader != nil {
		if existing.MasjidImageURL != "" {
			bucket, path, err := helper.ExtractSupabasePath(existing.MasjidImageURL)
			if err == nil {
				_ = helper.DeleteFromSupabase(bucket, path)
			}
		}
		url, err := helper.UploadImageAsWebPToSupabase("masjid", fileHeader)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal upload gambar"})
		}
		existing.MasjidImageURL = url
	} else {
		strURL := c.FormValue("masjid_image_url")
		if strURL != "" && strURL != existing.MasjidImageURL {
			if existing.MasjidImageURL != "" {
				bucket, path, err := helper.ExtractSupabasePath(existing.MasjidImageURL)
				if err == nil {
					_ = helper.DeleteFromSupabase(bucket, path)
				}
			}
			existing.MasjidImageURL = strURL
		}
	}

	// 🔁 Field lain (partial)
	formFields := map[string]*string{
		"masjid_name":          &existing.MasjidName,
		"masjid_bio_short":     &existing.MasjidBioShort,
		"masjid_location":      &existing.MasjidLocation,
		"masjid_instagram_url": &existing.MasjidInstagramURL,
		"masjid_whatsapp_url":  &existing.MasjidWhatsappURL,
		"masjid_youtube_url":   &existing.MasjidYoutubeURL,
	}
	for key, ptr := range formFields {
		val := c.FormValue(key)
		if val != "" {
			*ptr = val
		}
	}

	// Verifikasi
	if v := c.FormValue("masjid_is_verified"); v != "" {
		existing.MasjidIsVerified = v == "true"
	}

	// 💾 Simpan
	if err := mc.DB.Save(&existing).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui masjid"})
	}

	return c.JSON(fiber.Map{
		"message": "Masjid berhasil diperbarui",
		"data":    dto.FromModelMasjid(&existing),
	})
}

// 🟢 DELETE MASJID
func (mc *MasjidController) DeleteMasjid(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting masjid with ID: %s\n", id)

	if err := mc.DB.Delete(&model.MasjidModel{}, "masjid_id = ?", id).Error; err != nil {
		log.Printf("[ERROR] Failed to delete masjid: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghapus masjid",
		})
	}

	log.Printf("[SUCCESS] Masjid with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Masjid dengan ID %s berhasil dihapus", id),
	})
}
