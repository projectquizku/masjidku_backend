package dto

import (
	"masjidku_backend/internals/features/users/user/model"
	"time"

	"github.com/google/uuid"
)

// ✅ Request DTO untuk create/update profile
type UsersProfileRequest struct {
	DonationName string     `json:"donation_name"`
	FullName     string     `json:"full_name"`
	DateOfBirth  *time.Time `json:"date_of_birth"` // format: "2006-01-02"
	Gender       *string    `json:"gender"`        // "male" atau "female"
	PhoneNumber  string     `json:"phone_number"`
	Bio          string     `json:"bio"`
	Location     string     `json:"location"`
	Occupation   string     `json:"occupation"`
	ImageURL     *string    `json:"image_url"`
}

// ✅ Response DTO
type UsersProfileResponse struct {
	ID           uint      `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	DonationName string    `json:"donation_name"`
	FullName     string    `json:"full_name"`
	DateOfBirth  *string   `json:"date_of_birth"` // string agar aman saat null
	Gender       *string   `json:"gender,omitempty"`
	PhoneNumber  string    `json:"phone_number"`
	Bio          string    `json:"bio"`
	Location     string    `json:"location"`
	Occupation   string    `json:"occupation"`
	ImageURL     *string   `json:"image_url"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

// 🔄 Request → Model (untuk create/update)
func (r *UsersProfileRequest) ToModel(userID uuid.UUID) *model.UsersProfileModel {
	var gender *model.Gender
	if r.Gender != nil {
		g := model.Gender(*r.Gender)
		gender = &g
	}
	return &model.UsersProfileModel{
		UserID:       userID,
		DonationName: r.DonationName,
		FullName:     r.FullName,
		DateOfBirth:  r.DateOfBirth,
		Gender:       gender,
		PhoneNumber:  r.PhoneNumber,
		Bio:          r.Bio,
		Location:     r.Location,
		Occupation:   r.Occupation,
		ImageURL:     r.ImageURL,
	}
}

// 🔄 Model → Response
func ToUsersProfileResponse(m *model.UsersProfileModel) *UsersProfileResponse {
	var dob *string
	if m.DateOfBirth != nil {
		str := m.DateOfBirth.Format("2006-01-02")
		dob = &str
	}
	var gender *string
	if m.Gender != nil {
		g := string(*m.Gender)
		gender = &g
	}
	return &UsersProfileResponse{
		ID:           m.ID,
		UserID:       m.UserID,
		DonationName: m.DonationName,
		FullName:     m.FullName,
		DateOfBirth:  dob,
		Gender:       gender,
		PhoneNumber:  m.PhoneNumber,
		Bio:          m.Bio,
		Location:     m.Location,
		Occupation:   m.Occupation,
		ImageURL:     m.ImageURL,
		CreatedAt:    m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
