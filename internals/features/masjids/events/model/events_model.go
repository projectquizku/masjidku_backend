package model

import (
	"time"

	masjidmodel "masjidku_backend/internals/features/masjids/masjids/model"

	"github.com/google/uuid"
)

type EventModel struct {
	EventID                     uuid.UUID  `gorm:"column:event_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"event_id"`
	EventTitle                  string     `gorm:"column:event_title;type:varchar(255);not null" json:"event_title"`
	EventSlug                   string     `gorm:"column:event_slug;type:varchar(100);not null" json:"event_slug"`
	EventDescription            string     `gorm:"column:event_description;type:text" json:"event_description"`
	EventLocation               string     `gorm:"column:event_location;type:varchar(255)" json:"event_location"`
	EventImageURL               *string    `gorm:"column:event_image_url;type:text" json:"event_image_url"`
	EventCapacity               *int       `gorm:"column:event_capacity" json:"event_capacity"`
	EventIsPublic               bool       `gorm:"column:event_is_public;default:true" json:"event_is_public"`
	EventIsRegistrationRequired bool       `gorm:"column:event_is_registration_required;default:false" json:"event_is_registration_required"`
	EventIsPaid                 bool       `gorm:"column:event_is_paid;default:false" json:"event_is_paid"`
	EventPrice                  *int       `gorm:"column:event_price" json:"event_price"`
	EventPaymentDeadline        *time.Time `gorm:"column:event_payment_deadline" json:"event_payment_deadline"`

	EventMasjidID  uuid.UUID  `gorm:"column:event_masjid_id;type:uuid;not null" json:"event_masjid_id"`
	EventCreatedAt time.Time  `gorm:"column:event_created_at;autoCreateTime" json:"event_created_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	// Relasi
	Masjid *masjidmodel.MasjidModel `gorm:"foreignKey:EventMasjidID;references:MasjidID" json:"-"`
}

func (EventModel) TableName() string {
	return "events"
}
