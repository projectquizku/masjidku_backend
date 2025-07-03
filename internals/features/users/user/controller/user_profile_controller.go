package controller

import (
	"log"
	"time"

	"masjidku_backend/internals/features/users/user/dto"
	"masjidku_backend/internals/features/users/user/model"
	helper "masjidku_backend/internals/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersProfileController struct {
	DB *gorm.DB
}

func NewUsersProfileController(db *gorm.DB) *UsersProfileController {
	return &UsersProfileController{DB: db}
}

func (upc *UsersProfileController) GetProfiles(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all user profiles")

	var profiles []model.UsersProfileModel
	if err := upc.DB.Find(&profiles).Error; err != nil {
		log.Println("[ERROR] Failed to fetch user profiles:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to fetch user profiles")
	}

	// Konversi ke DTO response
	var responses []dto.UsersProfileResponse
	for _, p := range profiles {
		responses = append(responses, *dto.ToUsersProfileResponse(&p))
	}

	return helper.Success(c, "User profiles fetched successfully", responses)
}

func (upc *UsersProfileController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	log.Println("[INFO] Fetching user profile with user_id:", userID)

	var profile model.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	response := dto.ToUsersProfileResponse(&profile)
	return helper.Success(c, "User profile fetched successfully", response)
}

func (upc *UsersProfileController) CreateProfile(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user profile")

	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		log.Println("[ERROR] user_id not found or invalid in context")
		return helper.Error(c, fiber.StatusUnauthorized, "Unauthorized: no user_id")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", userIDStr)
		return helper.Error(c, fiber.StatusUnauthorized, "Unauthorized: invalid user_id")
	}

	// Ambil semua field dari multipart/form
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("[ERROR] Gagal membaca multipart form:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Form tidak valid (multipart/form-data diperlukan)")
	}
	get := func(key string) string {
		if form == nil || form.Value == nil || form.Value[key] == nil || len(form.Value[key]) == 0 {
			return ""
		}
		return form.Value[key][0]
	}

	// Parsing data form
	donationName := get("donation_name")
	fullName := get("full_name")
	dateStr := get("date_of_birth")
	genderStr := get("gender")
	phone := get("phone_number")
	bio := get("bio")
	location := get("location")
	occupation := get("occupation")

	// Parse date
	var dateOfBirth *time.Time
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			dateOfBirth = &parsed
		}
	}

	// Parse gender
	var gender *model.Gender
	if genderStr != "" {
		g := model.Gender(genderStr)
		gender = &g
	}

	// Proses image
	var imageURL *string
	if fileHeaders, ok := form.File["image_url"]; ok && len(fileHeaders) > 0 {
		file := fileHeaders[0]

		// Hapus gambar lama jika ada
		var old model.UsersProfileModel
		if err := upc.DB.Where("user_id = ?", userID).First(&old).Error; err == nil && old.ImageURL != nil {
			oldPath := helper.ExtractSupabaseStoragePath(*old.ImageURL)
			if oldPath != "" {
				_ = helper.DeleteFromSupabase("image", oldPath)
			}
		}

		url, err := helper.UploadImageAsWebPToSupabase("users/profile_images", file)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		imageURL = &url
	} else if val := get("image_url"); val != "" {
		imageURL = &val
	}

	// Siapkan model
	profile := model.UsersProfileModel{
		UserID:       userID,
		DonationName: donationName,
		FullName:     fullName,
		DateOfBirth:  dateOfBirth,
		Gender:       gender,
		PhoneNumber:  phone,
		Bio:          bio,
		Location:     location,
		Occupation:   occupation,
		ImageURL:     imageURL,
	}

	var existing model.UsersProfileModel
	tx := upc.DB.Where("user_id = ?", userID).First(&existing)

	if tx.RowsAffected > 0 {
		// Update profil lama
		if err := upc.DB.Model(&existing).Updates(profile).Error; err != nil {
			log.Println("[ERROR] Failed to update user profile:", err)
			return helper.Error(c, fiber.StatusInternalServerError, "Failed to update user profile")
		}
		log.Println("[SUCCESS] User profile updated:", userID)
		return helper.Success(c, "User profile updated successfully", dto.ToUsersProfileResponse(&existing))
	}

	// Create profil baru
	if err := upc.DB.Create(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to create user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to create user profile")
	}

	log.Println("[SUCCESS] User profile created:", userID)
	return helper.SuccessWithCode(c, fiber.StatusCreated, "User profile created successfully", dto.ToUsersProfileResponse(&profile))
}

func (upc *UsersProfileController) UpdateProfile(c *fiber.Ctx) error {
	log.Println("[INFO] Updating user profile")

	// ✅ Ambil user_id dari token (middleware)
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return helper.Error(c, fiber.StatusUnauthorized, "Unauthorized - user_id missing")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helper.Error(c, fiber.StatusUnauthorized, "Invalid user_id format")
	}

	// ✅ Cari profil
	var profile model.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	// ✅ Ambil data dari multipart/form
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("[ERROR] Gagal membaca multipart form:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Form tidak valid (multipart/form-data diperlukan)")
	}
	get := func(key string) string {
		if form == nil || form.Value == nil || form.Value[key] == nil || len(form.Value[key]) == 0 {
			return ""
		}
		return form.Value[key][0]
	}

	// ✅ Parsing form (semua opsional)
	if val := get("donation_name"); val != "" {
		profile.DonationName = val
	}
	if val := get("full_name"); val != "" {
		profile.FullName = val
	}
	if val := get("date_of_birth"); val != "" {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			profile.DateOfBirth = &parsed
		}
	}
	if val := get("gender"); val != "" {
		g := model.Gender(val)
		profile.Gender = &g
	}
	if val := get("phone_number"); val != "" {
		profile.PhoneNumber = val
	}
	if val := get("bio"); val != "" {
		profile.Bio = val
	}
	if val := get("location"); val != "" {
		profile.Location = val
	}
	if val := get("occupation"); val != "" {
		profile.Occupation = val
	}

	// ✅ Proses image (upload file atau string URL)
	if files, ok := form.File["image_url"]; ok && len(files) > 0 {
		file := files[0]

		// Hapus gambar lama jika ada
		if profile.ImageURL != nil {
			oldPath := helper.ExtractSupabaseStoragePath(*profile.ImageURL)
			_ = helper.DeleteFromSupabase("image", oldPath)
		}

		url, err := helper.UploadImageAsWebPToSupabase("users/profile_images", file)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		profile.ImageURL = &url
	} else if val := get("image_url"); val != "" {
		if profile.ImageURL != nil && *profile.ImageURL != val {
			oldPath := helper.ExtractSupabaseStoragePath(*profile.ImageURL)
			_ = helper.DeleteFromSupabase("image", oldPath)
		}
		profile.ImageURL = &val
	}

	// ✅ Simpan ke DB
	if err := upc.DB.Save(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to update user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to update user profile")
	}

	log.Println("[SUCCESS] User profile updated:", userID)
	return helper.Success(c, "User profile updated successfully", dto.ToUsersProfileResponse(&profile))
}

func (upc *UsersProfileController) DeleteProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	log.Println("[INFO] Deleting user profile with user_id:", userID)

	var profile model.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	if err := upc.DB.Delete(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to delete user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to delete user profile")
	}

	return helper.Success(c, "User profile deleted successfully", nil)
}
