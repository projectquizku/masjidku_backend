package routes

import (
	"log"
	"time"

	routeDetails "masjidku_backend/internals/route/details"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var startTime time.Time

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	startTime = time.Now()

	BaseRoutes(app, db)

	log.Println("[INFO] Setting up AuthRoutes...")
	routeDetails.AuthRoutes(app, db)

	log.Println("[INFO] Setting up UserRoutes...")
	routeDetails.UserRoutes(app, db)

	log.Println("[INFO] Setting up UtilsRoutes...")
	routeDetails.UtilsRoutes(app, db)

	log.Println("[INFO] Setting up CertificateRoutes...")
	routeDetails.CertificateRoutes(app, db)

	log.Println("[INFO] Setting up MasjidsRoutes")
	routeDetails.MasjidRoutes(app, db)

	log.Println("[INFO] Setting up HomeRoutes...")
	routeDetails.HomeRoutes(app, db)

}
