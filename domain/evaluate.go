package domain

import (
	"github.com/lib/pq"
)

type StudentEvaluation struct {
	ID                                int             `json:"id" gorm:"primaryKey autoIncrement"`
	StudentId                         string          `json:"studentId" gorm:"not null"`
	NewSources                        *pq.StringArray `json:"newSources" gorm:"type:text[]"`
	OverallActivity                   int             `json:"overallActivity"`
	InterestActivity                  int             `json:"interestActivity"`
	ReceivedFacultyInfoClearly        int             `json:"receivedFacultyInfoClearly"`
	WouldRecommendCUOpenHouseNextTime int             `json:"wouldRecommendCUOpenHouseNextTime"`
	FavoriteBooth                     *string         `json:"favoriteBooth"`
	ActivityDiversity                 int             `json:"activityDiversity"`
	PerceivedCrowdDensity             int             `json:"perceivedCrowdDensity"`
	HasFullBoothAccess                int             `json:"hasFullBoothAccess"`
	FacilityConvenienceRating         int             `json:"facilityConvenienceRating"`
	CampusNavigationRating            int             `json:"campusNavigationRating"`
	HesitationLevelAfterDisaster      int             `json:"hesitationLevelAfterDisaster"`
	LineOASignupRating                int             `json:"lineOASignupRating"`
	DesignBeautyRating                int             `json:"designBeautyRating"`
	WebsiteImprovementSuggestions     *string         `json:"websiteImprovementSuggestions"`

	Student User `gorm:"foreignKey:StudentId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
