package domain

import (
	"time"

	"github.com/lib/pq"
)

type Role string

const (
	Member  Role = "member"
	Student Role = "student"
	Staff   Role = "staff"
	Admin   Role = "admin"
)

type User struct {
	ID              string          `json:"id" gorm:"primaryKey"`
	UID             string          `json:"uid" gorm:"unique"`
	Name            string          `json:"name"`
	Role            Role            `json:"role"`
	Email           string          `json:"email"`
	Phone           string          `json:"phone" gorm:"unique"` // Make phone unique
	BirthDate       *time.Time      `json:"birthDate"`
	Status          *string         `json:"status"`      // ม.ต้น, ม.ปลาย, ปวช., ปวส. etc.
	OtherStatus     *string         `json:"otherStatus"` // other status
	Province        *string         `json:"province"`
	School          *string         `json:"school"`
	SelectedSources *pq.StringArray `json:"selectedSources" gorm:"type:text[]"`
	OtherSource     *string         `json:"otherSource"`
	FirstInterest   *string         `json:"firstInterest"`
	SecondInterest  *string         `json:"secondInterest"`
	ThirdInterest   *string         `json:"thirdInterest"`
	Objective       *string         `json:"objective"`
	RegisteredAt    *time.Time      `json:"registerAt"`
	LastEntered     *time.Time      `json:"lastEntered"` // Timestamp for the last QR scan

	// For staff/admin only
	Faculty        *string `json:"faculty"`
	StudentID      *string `json:"studentId"`
	Nickname       *string `json:"nickname"`
	Year           *int    `json:"year"`
	IsCentralStaff *bool   `json:"isCentralStaff"`
}

type StudentTransaction struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	StudentRegistrationID string    `json:"studentId"` // Foreign key index
	Faculty               string    `json:"faculty"`
	RegisteredAt          time.Time `json:"registeredAt"`

	// Relationship
	Student User `gorm:"foreignKey:StudentRegistrationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
