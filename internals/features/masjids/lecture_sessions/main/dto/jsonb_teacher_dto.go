package dto

import (
	"encoding/json"
	"fmt"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"
)

type JSONBTeacher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (j JSONBTeacher) ToModel() model.JSONBTeacher {
	return model.JSONBTeacher{
		ID:   j.ID,
		Name: j.Name,
	}
}

func FromModel(m model.JSONBTeacher) JSONBTeacher {
	return JSONBTeacher{
		ID:   m.ID,
		Name: m.Name,
	}
}

// 🔄 Konversi dari string JSON ke struct
func (t *JSONBTeacher) FromString(jsonStr string) error {
	if err := json.Unmarshal([]byte(jsonStr), t); err != nil {
		return fmt.Errorf("gagal parse JSONBTeacher: %w", err)
	}
	return nil
}
