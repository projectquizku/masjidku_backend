package route

import (
	"masjidku_backend/internals/features/masjids/lectures/events/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EventRoutesUser(api fiber.Router, db *gorm.DB) {
	// 🔹 Events (User hanya lihat, tidak bisa create)
	eventCtrl := controller.NewEventController(db)
	event := api.Group("/events")
	event.Get("/", eventCtrl.GetAllEvents)
	event.Post("/by-masjid", eventCtrl.GetEventsByMasjid)

	// 🔹 User Event Registrations
	registrationCtrl := controller.NewUserEventRegistrationController(db)
	reg := api.Group("/user-event-registrations")
	reg.Post("/", registrationCtrl.CreateRegistration)           // user daftar event
	reg.Post("/by-user", registrationCtrl.GetRegistrantsByEvent) // user lihat event yang diikuti
}
