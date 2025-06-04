package dto

import (
	"masjidku_backend/internals/features/masjids/lectures/events/model"
	"time"

	"github.com/google/uuid"
)

// Request untuk membuat event
type EventRequest struct {
	EventTitle                string    `json:"event_title"`
	EventDescription          string    `json:"event_description"`
	EventStartTime            time.Time `json:"event_start_time"`
	EventEndTime              time.Time `json:"event_end_time"`
	EventLocation             string    `json:"event_location"`
	EventIsRegistrationNeeded bool      `json:"event_is_registration_required"`
	EventCapacity             int       `json:"event_capacity"`
	EventImageURL             string    `json:"event_image_url"`
	EventMasjidID             uuid.UUID `json:"event_masjid_id"`
}

// Response untuk menampilkan event
type EventResponse struct {
	EventID                   uuid.UUID `json:"event_id"`
	EventTitle                string    `json:"event_title"`
	EventDescription          string    `json:"event_description"`
	EventStartTime            string    `json:"event_start_time"`
	EventEndTime              string    `json:"event_end_time"`
	EventLocation             string    `json:"event_location"`
	EventIsRegistrationNeeded bool      `json:"event_is_registration_required"`
	EventCapacity             int       `json:"event_capacity"`
	EventImageURL             string    `json:"event_image_url"`
	EventMasjidID             uuid.UUID `json:"event_masjid_id"`
	EventCreatedAt            string    `json:"event_created_at"`
}

// Konversi dari request → model
func (r *EventRequest) ToModel() *model.EventModel {
	return &model.EventModel{
		EventTitle:                r.EventTitle,
		EventDescription:          r.EventDescription,
		EventStartTime:            r.EventStartTime,
		EventEndTime:              r.EventEndTime,
		EventLocation:             r.EventLocation,
		EventIsRegistrationNeeded: r.EventIsRegistrationNeeded,
		EventCapacity:             r.EventCapacity,
		EventImageURL:             r.EventImageURL,
		EventMasjidID:             r.EventMasjidID,
	}
}

// Konversi dari model → response
func ToEventResponse(m *model.EventModel) *EventResponse {
	return &EventResponse{
		EventID:                   m.EventID,
		EventTitle:                m.EventTitle,
		EventDescription:          m.EventDescription,
		EventStartTime:            m.EventStartTime.Format("2006-01-02 15:04:05"),
		EventEndTime:              m.EventEndTime.Format("2006-01-02 15:04:05"),
		EventLocation:             m.EventLocation,
		EventIsRegistrationNeeded: m.EventIsRegistrationNeeded,
		EventCapacity:             m.EventCapacity,
		EventImageURL:             m.EventImageURL,
		EventMasjidID:             m.EventMasjidID,
		EventCreatedAt:            m.EventCreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Konversi list model → list response
func ToEventResponseList(models []model.EventModel) []EventResponse {
	var result []EventResponse
	for _, m := range models {
		result = append(result, *ToEventResponse(&m))
	}
	return result
}
