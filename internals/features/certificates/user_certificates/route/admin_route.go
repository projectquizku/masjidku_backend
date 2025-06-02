// routes/issued_certificate_admin_routes.go
package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	controller "masjidku_backend/internals/features/certificates/user_certificates/controller"
)

func IssuedCertificateAdminRoutes(app fiber.Router, db *gorm.DB) {
	c := controller.NewIssuedCertificateController(db)
	app.Get("/:id", c.GetByIDUser) // GET issued certificate by ID (admin only)
}
