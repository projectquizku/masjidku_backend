package dto

import (
	"masjidku_backend/internals/features/masjids/events/model"
	"time"

	"github.com/google/uuid"
)

// ✅ Request DTO (disesuaikan dengan SQL)
type EventSessionRequest struct {
	EventSessionEventID     uuid.UUID  `json:"event_session_event_id"`
	EventSessionTitle       string     `json:"event_session_title"`
	EventSessionDescription string     `json:"event_session_description"`
	EventSessionStartTime   time.Time  `json:"event_session_start_time"`
	EventSessionEndTime     time.Time  `json:"event_session_end_time"`
	EventSessionLocation    string     `json:"event_session_location"`
	EventSessionImageURL    string     `json:"event_session_image_url"`
	EventSessionMasjidID    uuid.UUID  `json:"event_session_masjid_id"`
	EventSessionCreatedBy   *uuid.UUID `json:"event_session_created_by"` // bisa null
}

// ✅ Response DTO
type EventSessionResponse struct {
	EventSessionID          uuid.UUID  `json:"event_session_id"`
	EventSessionEventID     uuid.UUID  `json:"event_session_event_id"`
	EventSessionTitle       string     `json:"event_session_title"`
	EventSessionDescription string     `json:"event_session_description"`
	EventSessionStartTime   string     `json:"event_session_start_time"`
	EventSessionEndTime     string     `json:"event_session_end_time"`
	EventSessionLocation    string     `json:"event_session_location"`
	EventSessionImageURL    string     `json:"event_session_image_url"`
	EventSessionMasjidID    uuid.UUID  `json:"event_session_masjid_id"`
	EventSessionCreatedBy   *uuid.UUID `json:"event_session_created_by"`
	EventSessionCreatedAt   string     `json:"event_session_created_at"`
	EventSessionUpdatedAt   string     `json:"event_session_updated_at"`
}

// 🔄 Request → Model
func (r *EventSessionRequest) ToModel() *model.EventSessionModel {
	return &model.EventSessionModel{
		EventSessionEventID:     r.EventSessionEventID,
		EventSessionTitle:       r.EventSessionTitle,
		EventSessionDescription: r.EventSessionDescription,
		EventSessionStartTime:   r.EventSessionStartTime,
		EventSessionEndTime:     r.EventSessionEndTime,
		EventSessionLocation:    r.EventSessionLocation,
		EventSessionImageURL:    r.EventSessionImageURL,
		EventSessionMasjidID:    r.EventSessionMasjidID,
		EventSessionCreatedBy:   r.EventSessionCreatedBy,
	}
}

// 🔄 Model → Response
func ToEventSessionResponse(m *model.EventSessionModel) *EventSessionResponse {
	return &EventSessionResponse{
		EventSessionID:          m.EventSessionID,
		EventSessionEventID:     m.EventSessionEventID,
		EventSessionTitle:       m.EventSessionTitle,
		EventSessionDescription: m.EventSessionDescription,
		EventSessionStartTime:   m.EventSessionStartTime.Format("2006-01-02 15:04:05"),
		EventSessionEndTime:     m.EventSessionEndTime.Format("2006-01-02 15:04:05"),
		EventSessionLocation:    m.EventSessionLocation,
		EventSessionImageURL:    m.EventSessionImageURL,
		EventSessionMasjidID:    m.EventSessionMasjidID,
		EventSessionCreatedBy:   m.EventSessionCreatedBy,
		EventSessionCreatedAt:   m.EventSessionCreatedAt.Format("2006-01-02 15:04:05"),
		EventSessionUpdatedAt:   m.EventSessionUpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// 🔄 List Model → List Response
func ToEventSessionResponseList(models []model.EventSessionModel) []EventSessionResponse {
	var responses []EventSessionResponse
	for _, m := range models {
		responses = append(responses, *ToEventSessionResponse(&m))
	}
	return responses
}
