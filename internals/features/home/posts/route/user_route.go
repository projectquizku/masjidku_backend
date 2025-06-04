package route

import (
	"masjidku_backend/internals/features/home/posts/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PostUserRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewPostController(db)

	user := api.Group("/posts")

	user.Get("/", ctrl.GetAllPosts)    // 📄 Semua post publik
	user.Get("/:id", ctrl.GetPostByID) // 🔍 Detail post
	// (opsional: bisa tambahkan route untuk like/unlike post di sini nanti)

	ctrl2 := controller.NewPostLikeController(db)

	post := api.Group("/post-likes")

	// 🔄 Toggle like (user harus login → ambil user_id dari token)
	post.Post("/toggle", ctrl2.ToggleLike)
}
