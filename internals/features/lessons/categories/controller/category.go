package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	dto "masjidku_backend/internals/features/lessons/categories/dto"
	"masjidku_backend/internals/features/lessons/categories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

// 🟢 GET ALL CATEGORIES: Ambil semua kategori dari database
func (cc *CategoryController) GetCategories(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all categories")

	var categories []model.CategoryModel

	// 🔍 Ambil semua kategori beserta relasi subcategories
	if err := cc.DB.Preload("Subcategories").Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data kategori",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d categories\n", len(categories))
	return c.JSON(fiber.Map{
		"message": "Data semua kategori berhasil diambil",
		"total":   len(categories),
		"data":    categories,
	})
}

// 🟢 GET CATEGORY BY ID: Ambil satu kategori berdasarkan ID
func (cc *CategoryController) GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching category with ID: %s\n", id)

	var category model.CategoryModel

	// 🔍 Cari berdasarkan ID dan preload relasi
	if err := cc.DB.Preload("Subcategories").
		Where("category_id = ?", id).
		First(&category).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{
			"error": "Kategori tidak ditemukan",
		})
	}

	log.Printf("[SUCCESS] Retrieved category: ID=%s, Name=%s\n", id, category.CategoryName)
	return c.JSON(fiber.Map{
		"message": "Data kategori berhasil diambil",
		"data":    category,
	})
}

// 🟢 GET CATEGORY NAMES BY DIFFICULTY ID: Ambil nama kategori berdasarkan difficulty_id
func (cc *CategoryController) GetCategoriesByDifficulty(c *fiber.Ctx) error {
	// 🔧 Ubah dari string ke integer, biar aman jika kolom bertipe int di DB
	difficultyIDParam := c.Params("difficulty_id")
	difficultyID, err := strconv.Atoi(difficultyIDParam)
	if err != nil {
		log.Printf("[ERROR] Invalid difficulty_id param: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Parameter difficulty_id tidak valid",
		})
	}

	log.Printf("[INFO] Fetching categories with difficulty ID: %d\n", difficultyID)

	var categories []model.CategoryModel

	// ✅ Revisi nama kolom jika ternyata di Supabase bukan 'category_difficulty_id' tapi 'difficulty_id'
	if err := cc.DB.
		Select("category_id", "category_name").
		Where("category_difficulty_id = ?", difficultyID). // ← ganti jika di Supabase memang "difficulty_id"
		Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories for difficulty ID %d: %v\n", difficultyID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil data kategori",
			"detail":  err.Error(), // untuk debugging, boleh dihapus di production
		})
	}

	// 🔄 Format response ke dalam DTO
	var responses []dto.CategoryTooltipResponse
	for _, cat := range categories {
		responses = append(responses, dto.CategoryTooltipResponse{
			CategoriesID:   cat.CategoryID,
			CategoriesName: cat.CategoryName,
		})
	}

	log.Printf("[SUCCESS] Retrieved %d category tooltips for difficulty ID %d\n", len(responses), difficultyID)
	return c.JSON(fiber.Map{
		"message": "Nama-nama kategori berhasil diambil",
		"total":   len(responses),
		"data":    responses,
	})
}

// 🟢 CREATE CATEGORY: Tambahkan satu atau banyak kategori
func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create category")

	var single model.CategoryModel
	var multiple []model.CategoryModel

	// 🌀 Parsing array jika dikirim banyak data
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		if err := cc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to create multiple categories: %v\n", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal menyimpan banyak kategori",
			})
		}
		log.Printf("[SUCCESS] %d categories created\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Kategori berhasil dibuat (multiple)",
			"data":    multiple,
		})
	}

	// 🌀 Jika bukan array, parsing objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Invalid input for single category: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Format input tidak valid",
		})
	}

	if err := cc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to create single category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menyimpan kategori",
		})
	}

	log.Printf("[SUCCESS] Category created: ID=%d, Name=%s\n", single.CategoryID, single.CategoryName)
	return c.Status(201).JSON(fiber.Map{
		"message": "Kategori berhasil dibuat",
		"data":    single,
	})
}

// 🟢 UPDATE CATEGORY: Perbarui kategori berdasarkan ID
func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating category with ID: %s\n", id)

	var category model.CategoryModel
	if err := cc.DB.Where("categories_id = ?", id).First(&category).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{
			"error": "Kategori tidak ditemukan",
		})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Invalid input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Input tidak valid",
		})
	}

	if raw, ok := input["categories_update_news"]; ok {
		if jsonData, err := json.Marshal(raw); err == nil {
			input["categories_update_news"] = datatypes.JSON(jsonData)
		}
	}

	if err := cc.DB.Model(&category).Updates(input).Error; err != nil {
		log.Printf("[ERROR] Failed to update category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal memperbarui kategori",
		})
	}

	log.Printf("[SUCCESS] Category updated: ID=%s, Name=%s\n", id, category.CategoryName)
	return c.JSON(fiber.Map{
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

// 🟢 DELETE CATEGORY: Hapus kategori berdasarkan ID
func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting category with ID: %s\n", id)

	if err := cc.DB.
		Where("categories_id = ?", id).
		Delete(&model.CategoryModel{}).Error; err != nil {
		log.Printf("[ERROR] Failed to delete category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghapus kategori",
		})
	}

	log.Printf("[SUCCESS] Category with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Kategori dengan ID %s berhasil dihapus", id),
	})
}
