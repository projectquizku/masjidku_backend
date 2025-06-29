package dto

import (
	"masjidku_backend/internals/features/home/posts/model"
	"time"
)

// ============================
// Response DTO
// ============================
type PostDTO struct {
	PostID          string     `json:"post_id"`
	PostTitle       string     `json:"post_title"`
	PostContent     string     `json:"post_content"`
	PostImageURL    *string    `json:"post_image_url"`
	PostIsPublished bool       `json:"post_is_published"`
	PostType        string     `json:"post_type"`
	PostMasjidID    *string    `json:"post_masjid_id"`
	PostUserID      *string    `json:"post_user_id"`
	PostCreatedAt   time.Time  `json:"post_created_at"`
	PostUpdatedAt   time.Time  `json:"post_updated_at"`
	PostDeletedAt   *time.Time `json:"post_deleted_at"`
}

// ============================
// Create Request DTO
// ============================
type CreatePostRequest struct {
	PostTitle       string  `json:"post_title" validate:"required,min=3"`
	PostContent     string  `json:"post_content" validate:"required"`
	PostImageURL    *string `json:"post_image_url"`
	PostIsPublished bool    `json:"post_is_published"`
	PostType        string  `json:"post_type" validate:"omitempty,oneof=text image video"`
	PostMasjidID    *string `json:"post_masjid_id"`
}

// ============================
// Update Request DTO (semua opsional)
// ============================
type UpdatePostRequest struct {
	PostTitle       *string `json:"post_title"`
	PostContent     *string `json:"post_content"`
	PostImageURL    *string `json:"post_image_url"` // boleh dikosongkan
	PostIsPublished *bool   `json:"post_is_published"`
	PostType        *string `json:"post_type"`
}

// ============================
// Converter
// ============================
func ToPostDTO(m model.PostModel) PostDTO {
	return PostDTO{
		PostID:          m.PostID,
		PostTitle:       m.PostTitle,
		PostContent:     m.PostContent,
		PostImageURL:    m.PostImageURL,
		PostIsPublished: m.PostIsPublished,
		PostType:        m.PostType,
		PostMasjidID:    m.PostMasjidID,
		PostUserID:      m.PostUserID,
		PostCreatedAt:   m.PostCreatedAt,
		PostUpdatedAt:   m.PostUpdatedAt,
		PostDeletedAt:   m.PostDeletedAt,
	}
}

func ToPostModel(req CreatePostRequest, userID *string) model.PostModel {
	return model.PostModel{
		PostTitle:       req.PostTitle,
		PostContent:     req.PostContent,
		PostImageURL:    req.PostImageURL,
		PostIsPublished: req.PostIsPublished,
		PostType:        req.PostType,
		PostMasjidID:    req.PostMasjidID,
		PostUserID:      userID,
	}
}

func UpdatePostModel(m *model.PostModel, req UpdatePostRequest) {
	if req.PostTitle != nil {
		m.PostTitle = *req.PostTitle
	}
	if req.PostContent != nil {
		m.PostContent = *req.PostContent
	}
	if req.PostImageURL != nil {
		m.PostImageURL = req.PostImageURL
	}
	if req.PostIsPublished != nil {
		m.PostIsPublished = *req.PostIsPublished
	}
	if req.PostType != nil {
		m.PostType = *req.PostType
	}
}
