package model

import (
	"time"

	"github.com/google/uuid"
)

type EventModel struct {
	EventID                   uuid.UUID `gorm:"column:event_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"event_id"`
	EventTitle                string    `gorm:"column:event_title;type:varchar(255);not null" json:"event_title"`
	EventDescription          string    `gorm:"column:event_description;type:text" json:"event_description"`
	EventStartTime            time.Time `gorm:"column:event_start_time;not null" json:"event_start_time"`
	EventEndTime              time.Time `gorm:"column:event_end_time;not null" json:"event_end_time"`
	EventLocation             string    `gorm:"column:event_location;type:varchar(255)" json:"event_location"`
	EventIsRegistrationNeeded bool      `gorm:"column:event_is_registration_required;default:false" json:"event_is_registration_required"`
	EventCapacity             int       `gorm:"column:event_capacity" json:"event_capacity"`
	EventImageURL             string    `gorm:"column:event_image_url;type:text" json:"event_image_url"`
	EventCreatedAt            time.Time `gorm:"column:event_created_at;autoCreateTime" json:"event_created_at"`
	EventMasjidID             uuid.UUID `gorm:"column:event_masjid_id;type:uuid;not null" json:"event_masjid_id"`
}

// TableName overrides the default table name
func (EventModel) TableName() string {
	return "events"
}
