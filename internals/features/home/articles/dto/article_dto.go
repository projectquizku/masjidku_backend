package dto

import (
	"masjidku_backend/internals/features/home/articles/model"
	"time"
)

// ============================
// Response DTO
// ============================

type ArticleDTO struct {
	ArticleID          string    `json:"article_id"`
	ArticleTitle       string    `json:"article_title"`
	ArticleDescription string    `json:"article_description"`
	ArticleImageURL    string    `json:"article_image_url"`
	ArticleOrderID     int       `json:"article_order_id"`
	ArticleCreatedAt   time.Time `json:"article_created_at"`
	ArticleUpdatedAt   time.Time `json:"article_updated_at"`
}

// ============================
// Create & Update Request DTO
// ============================

type CreateArticleRequest struct {
	ArticleTitle       string `json:"article_title" form:"article_title" validate:"required,min=3"`
	ArticleDescription string `json:"article_description" form:"article_description" validate:"required"`
	ArticleImageURL    string `json:"article_image_url" form:"article_image_url"`
	ArticleOrderID     int    `json:"article_order_id" form:"article_order_id"`
}

type UpdateArticleRequest struct {
	ArticleTitle       string `json:"article_title" form:"article_title" validate:"omitempty,min=3"`
	ArticleDescription string `json:"article_description" form:"article_description" validate:"omitempty"`
	ArticleImageURL    string `json:"article_image_url" form:"article_image_url"`
	ArticleOrderID     *int   `json:"article_order_id" form:"article_order_id"` // pointer agar bisa deteksi tidak dikirim
}

// ============================
// Converter
// ============================

func ToArticleDTO(m model.ArticleModel) ArticleDTO {
	return ArticleDTO{
		ArticleID:          m.ArticleID,
		ArticleTitle:       m.ArticleTitle,
		ArticleDescription: m.ArticleDescription,
		ArticleImageURL:    m.ArticleImageURL,
		ArticleOrderID:     m.ArticleOrderID,
		ArticleCreatedAt:   m.ArticleCreatedAt,
		ArticleUpdatedAt:   m.ArticleUpdatedAt,
	}
}
