package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Struktur JSON untuk kolom lecture_session_teacher (JSONB)
type JSONBTeacher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (j JSONBTeacher) Value() (driver.Value, error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (j *JSONBTeacher) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONBTeacher value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

type LectureSessionModel struct {
	LectureSessionID            uuid.UUID    `gorm:"column:lecture_session_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_session_id"`
	LectureSessionTitle         string       `gorm:"column:lecture_session_title;type:varchar(255);not null" json:"lecture_session_title"`
	LectureSessionDescription   string       `gorm:"column:lecture_session_description;type:text" json:"lecture_session_description"`
	LectureSessionTeacher       JSONBTeacher `gorm:"column:lecture_session_teacher;type:jsonb;not null" json:"lecture_session_teacher"`
	LectureSessionStartTime     time.Time    `gorm:"column:lecture_session_start_time;not null" json:"lecture_session_start_time"`
	LectureSessionEndTime       time.Time    `gorm:"column:lecture_session_end_time;not null" json:"lecture_session_end_time"`
	LectureSessionPlace         *string      `gorm:"column:lecture_session_place;type:text" json:"lecture_session_place"`
	LectureSessionImageURL      *string      `gorm:"column:lecture_session_image_url;type:text" json:"lecture_session_image_url"`
	LectureSessionLectureID     *uuid.UUID   `gorm:"column:lecture_session_lecture_id;type:uuid" json:"lecture_session_lecture_id"`
	LectureSessionCertificateID *uuid.UUID   `gorm:"column:lecture_session_certificate_id;type:uuid" json:"lecture_session_certificate_id,omitempty"`
	LectureSessionCreatedAt     time.Time    `gorm:"column:lecture_session_created_at;autoCreateTime" json:"lecture_session_created_at"`
}

func (LectureSessionModel) TableName() string {
	return "lecture_sessions"
}

func SyncTotalLectureSessions(db *gorm.DB, lectureID uuid.UUID) error {
	log.Println("[SERVICE] SyncTotalLectureSessions - lectureID:", lectureID)

	return db.Exec(`
		UPDATE lectures
		SET total_lecture_sessions = (
			SELECT COUNT(*) FROM lecture_sessions
			WHERE lecture_session_lecture_id = ?
		)
		WHERE lecture_id = ?
	`, lectureID, lectureID).Error
}

func (s *LectureSessionModel) AfterSave(tx *gorm.DB) error {
	if s.LectureSessionLectureID != nil {
		return SyncTotalLectureSessions(tx, *s.LectureSessionLectureID)
	}
	return nil
}

func (s *LectureSessionModel) AfterDelete(tx *gorm.DB) error {
	log.Printf("[HOOK] AfterDelete triggered for LectureSessionID: %s", s.LectureSessionID)

	var lectureID uuid.UUID
	if err := tx.Unscoped().
		Model(&LectureSessionModel{}).
		Select("lecture_session_lecture_id").
		Where("lecture_session_id = ?", s.LectureSessionID).
		Take(&lectureID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch lecture_session_lecture_id after delete:", err)
		return err
	}

	log.Printf("[HOOK] Fetched lecture_session_lecture_id: %s for deleted session", lectureID)
	return SyncTotalLectureSessions(tx, lectureID)
}
