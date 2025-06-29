package controller

import (
	"fmt"
	"masjidku_backend/internals/features/home/posts/dto"
	"masjidku_backend/internals/features/home/posts/model"
	helper "masjidku_backend/internals/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validatePost = validator.New()

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

func (ctrl *PostController) CreatePost(c *fiber.Ctx) error {
	// 🔍 Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid multipart/form-data")
	}

	// 🎯 Ambil field teks dari form
	postTitle := c.FormValue("post_title")
	postContent := c.FormValue("post_content")
	postType := c.FormValue("post_type")
	postIsPublished := c.FormValue("post_is_published") == "true"
	postMasjidID := c.FormValue("post_masjid_id")

	// 🆔 Ambil user ID dari token
	userID := c.Locals("user_id").(string)

	// 🖼️ Proses upload gambar dari file ATAU gunakan URL langsung
	var imageURL *string
	files := form.File["post_image_url"]
	if len(files) > 0 {
		uploaded := files[0]
		url, err := helper.UploadImageAsWebPToSupabase("posts", uploaded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Gagal upload gambar: %v", err))
		}
		imageURL = &url
	} else {
		strURL := c.FormValue("post_image_url")
		if strURL != "" {
			imageURL = &strURL
		}
	}

	// 🧱 Buat model post
	post := model.PostModel{
		PostTitle:       postTitle,
		PostContent:     postContent,
		PostImageURL:    imageURL,
		PostIsPublished: postIsPublished,
		PostType:        postType,
		PostMasjidID:    &postMasjidID,
		PostUserID:      &userID,
	}

	// 💾 Simpan ke database
	if err := ctrl.DB.Create(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan postingan")
	}

	// 📨 Return response ke frontend
	return c.Status(fiber.StatusCreated).JSON(dto.ToPostDTO(post))
}

// 🔄 Update Post
func (ctrl *PostController) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	var post model.PostModel
	if err := ctrl.DB.First(&post, "post_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "post not found")
	}

	// 🖼️ Proses gambar: file atau string
	var newImageURL *string
	fileHeader, err := c.FormFile("post_image_url")
	if err == nil && fileHeader != nil {
		if post.PostImageURL != nil {
			if bucket, path, err := helper.ExtractSupabasePath(*post.PostImageURL); err == nil {
				_ = helper.DeleteFromSupabase(bucket, path)
			}
		}
		url, err := helper.UploadImageAsWebPToSupabase("posts", fileHeader)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "gagal upload gambar")
		}
		newImageURL = &url
	} else {
		urlStr := c.FormValue("post_image_url")
		if urlStr != "" {
			if post.PostImageURL != nil && *post.PostImageURL != urlStr {
				if bucket, path, err := helper.ExtractSupabasePath(*post.PostImageURL); err == nil {
					_ = helper.DeleteFromSupabase(bucket, path)
				}
			}
			newImageURL = &urlStr
		}
	}

	// 🔄 Update hanya field yang dikirim
	if req.PostTitle != nil {
		post.PostTitle = *req.PostTitle
	}
	if req.PostContent != nil {
		post.PostContent = *req.PostContent
	}
	if req.PostIsPublished != nil {
		post.PostIsPublished = *req.PostIsPublished
	}
	if req.PostType != nil {
		post.PostType = *req.PostType
	}
	if newImageURL != nil {
		post.PostImageURL = newImageURL
	}

	// ✅ Simpan
	if err := ctrl.DB.Save(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update post")
	}

	return c.JSON(dto.ToPostDTO(post))
}

// 📄 Get Semua Post
func (ctrl *PostController) GetAllPosts(c *fiber.Ctx) error {
	var posts []model.PostModel
	if err := ctrl.DB.Preload("Masjid").Preload("User").Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve posts")
	}

	var result []dto.PostDTO
	for _, post := range posts {
		result = append(result, dto.ToPostDTO(post))
	}

	return c.JSON(result)
}

// 🔍 Get Post by ID
func (ctrl *PostController) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var post model.PostModel
	if err := ctrl.DB.Preload("Masjid").Preload("User").First(&post, "post_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	return c.JSON(dto.ToPostDTO(post))
}

// 🗑️ Hapus Post
func (ctrl *PostController) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.PostModel{}, "post_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete post")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// =============================
// 📄 Get Posts by Masjid ID
// =============================
func (ctrl *PostController) GetPostsByMasjid(c *fiber.Ctx) error {
	type RequestBody struct {
		MasjidID string `json:"masjid_id" validate:"required,uuid"`
	}

	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// ✅ Ganti validate → validatePost
	if err := validatePost.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var posts []model.PostModel
	if err := ctrl.DB.
		Where("post_masjid_id = ? AND post_deleted_at IS NULL", req.MasjidID).
		Order("post_created_at DESC").
		Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve posts")
	}

	var result []dto.PostDTO
	for _, post := range posts {
		result = append(result, dto.ToPostDTO(post))
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil daftar postingan masjid",
		"data":    result,
	})
}
