package domain

type StudentEvaluation struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	StudentRegistrationID string    `json:"studentId"` // Foreign key index

	// Relationship
	Student User `gorm:"foreignKey:StudentRegistrationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}