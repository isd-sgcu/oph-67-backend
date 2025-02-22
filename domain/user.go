package domain

import "time"

type Role string

const (
	Member Role = "member"
	Staff  Role = "staff"
	Admin  Role = "admin"
)

type User struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	UID             string     `json:"uid" gorm:"unique"`
	Name            string     `json:"name"`
	Email           *string    `json:"email"`
	Phone           string     `json:"phone" gorm:"unique"` // Make phone unique
	BirtDate        *time.Time `json:"birthDate"`
	Role            Role       `json:"role"`
	Province        string     `json:"province"`
	School          string     `json:"school"`
	SelectedSources []string   `json:"selectedSources" gorm:"type:text[]"`
	OtherSource     *string    `json:"otherSource"`
	FirstInterest   string     `json:"firstInterest"`
	SecondInterest  string     `json:"secondInterest"`
	ThirdInterest   string     `json:"thirdInterest"`
	Objective       string     `json:"objective"`
	RegisteredAt    *time.Time `json:"registerAt"`
	LastEntered     *time.Time `json:"lastEntered"` // Timestamp for the last QR scan
}
