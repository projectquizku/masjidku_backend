package controller

import (
	"masjidku_backend/internals/features/home/qoutes/dto"
	"masjidku_backend/internals/features/home/qoutes/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validateQuote = validator.New()

type QuoteController struct {
	DB *gorm.DB
}

func NewQuoteController(db *gorm.DB) *QuoteController {
	return &QuoteController{DB: db}
}

// =============================
// ➕ Create Quote
// =============================
func (ctrl *QuoteController) CreateQuote(c *fiber.Ctx) error {
	var body dto.CreateQuoteRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := validateQuote.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	quote := model.QuoteModel{
		QuoteText:    body.QuoteText,
		IsPublished:  body.IsPublished,
		DisplayOrder: body.DisplayOrder,
	}

	if err := ctrl.DB.Create(&quote).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create quote")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToQuoteDTO(quote))
}

// =============================
// 🔄 Update Quote
// =============================
func (ctrl *QuoteController) UpdateQuote(c *fiber.Ctx) error {
	id := c.Params("id")

	var body dto.UpdateQuoteRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := validateQuote.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var quote model.QuoteModel
	if err := ctrl.DB.First(&quote, "quote_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Quote not found")
	}

	quote.QuoteText = body.QuoteText
	quote.IsPublished = body.IsPublished
	quote.DisplayOrder = body.DisplayOrder
	quote.CreatedAt = time.Now()

	if err := ctrl.DB.Save(&quote).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update quote")
	}

	return c.JSON(dto.ToQuoteDTO(quote))
}

// =============================
// 🗑️ Delete Quote
// =============================
func (ctrl *QuoteController) DeleteQuote(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.QuoteModel{}, "quote_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete quote")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// =============================
// 📄 Get All Quotes
// =============================
func (ctrl *QuoteController) GetAllQuotes(c *fiber.Ctx) error {
	var quotes []model.QuoteModel
	if err := ctrl.DB.Order("display_order ASC").Find(&quotes).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve quotes")
	}

	var result []dto.QuoteDTO
	for _, q := range quotes {
		result = append(result, dto.ToQuoteDTO(q))
	}

	return c.JSON(result)
}

// =============================
// 🔍 Get Quote By ID
// =============================
func (ctrl *QuoteController) GetQuoteByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var quote model.QuoteModel
	if err := ctrl.DB.First(&quote, "quote_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Quote not found")
	}

	return c.JSON(dto.ToQuoteDTO(quote))
}
