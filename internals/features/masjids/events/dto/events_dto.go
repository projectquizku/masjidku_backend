package dto

import (
	"strings"
	"time"

	"masjidku_backend/internals/features/masjids/events/model"

	"github.com/google/uuid"
)

// 🔹 Request untuk membuat event
type EventRequest struct {
	EventTitle                  string     `json:"event_title"`
	EventDescription            string     `json:"event_description"`
	EventLocation               string     `json:"event_location"`
	EventImageURL               *string    `json:"event_image_url"`
	EventCapacity               *int       `json:"event_capacity"`
	EventIsPublic               bool       `json:"event_is_public"`
	EventIsRegistrationRequired bool       `json:"event_is_registration_required"`
	EventIsPaid                 bool       `json:"event_is_paid"`
	EventPrice                  *int       `json:"event_price"`
	EventPaymentDeadline        *time.Time `json:"event_payment_deadline"`
	EventMasjidID               uuid.UUID  `json:"event_masjid_id"`
}

// 🔹 Response untuk menampilkan event
type EventResponse struct {
	EventID                     uuid.UUID  `json:"event_id"`
	EventTitle                  string     `json:"event_title"`
	EventSlug                   string     `json:"event_slug"`
	EventDescription            string     `json:"event_description"`
	EventLocation               string     `json:"event_location"`
	EventImageURL               *string    `json:"event_image_url"`
	EventCapacity               *int       `json:"event_capacity"`
	EventIsPublic               bool       `json:"event_is_public"`
	EventIsRegistrationRequired bool       `json:"event_is_registration_required"`
	EventIsPaid                 bool       `json:"event_is_paid"`
	EventPrice                  *int       `json:"event_price"`
	EventPaymentDeadline        *time.Time `json:"event_payment_deadline"`
	EventMasjidID               uuid.UUID  `json:"event_masjid_id"`
	EventCreatedAt              string     `json:"event_created_at"`
}

// 🔄 Fungsi bantu generate slug dari judul
func generateSlug(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

// 🔄 Konversi dari request → model
func (r *EventRequest) ToModel() *model.EventModel {
	return &model.EventModel{
		EventTitle:                  r.EventTitle,
		EventSlug:                   generateSlug(r.EventTitle),
		EventDescription:            r.EventDescription,
		EventLocation:               r.EventLocation,
		EventImageURL:               r.EventImageURL,
		EventCapacity:               r.EventCapacity,
		EventIsPublic:               r.EventIsPublic,
		EventIsRegistrationRequired: r.EventIsRegistrationRequired,
		EventIsPaid:                 r.EventIsPaid,
		EventPrice:                  r.EventPrice,
		EventPaymentDeadline:        r.EventPaymentDeadline,
		EventMasjidID:               r.EventMasjidID,
	}
}

// 🔄 Konversi dari model → response
func ToEventResponse(m *model.EventModel) *EventResponse {
	return &EventResponse{
		EventID:                     m.EventID,
		EventTitle:                  m.EventTitle,
		EventSlug:                   m.EventSlug,
		EventDescription:            m.EventDescription,
		EventLocation:               m.EventLocation,
		EventImageURL:               m.EventImageURL,
		EventCapacity:               m.EventCapacity,
		EventIsPublic:               m.EventIsPublic,
		EventIsRegistrationRequired: m.EventIsRegistrationRequired,
		EventIsPaid:                 m.EventIsPaid,
		EventPrice:                  m.EventPrice,
		EventPaymentDeadline:        m.EventPaymentDeadline,
		EventMasjidID:               m.EventMasjidID,
		EventCreatedAt:              m.EventCreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// 🔄 Konversi list model → list response
func ToEventResponseList(models []model.EventModel) []EventResponse {
	var result []EventResponse
	for _, m := range models {
		result = append(result, *ToEventResponse(&m))
	}
	return result
}
