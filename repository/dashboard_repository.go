package repository

import (
	"github.com/isd-sgcu/oph-67-backend/domain"
	"gorm.io/gorm"
)

type DashBoardRepository struct {
	DB *gorm.DB
}

func NewDashBoardRepository(db *gorm.DB) *DashBoardRepository {
	return &DashBoardRepository{DB: db}
}

func (r *DashBoardRepository) GetFacultyCount() ([]domain.FacultyPercent, error) {
	var results []domain.FacultyPercent

	query := `
        SELECT
            t.faculty,
            SUM(t.first_count) AS first_interest,
            SUM(t.second_count) AS second_interest,
            SUM(t.third_count) AS third_interest
        FROM users
        CROSS JOIN LATERAL (
            VALUES 
                (users.first_interest, 1, 0, 0),
                (users.second_interest, 0, 1, 0),
                (users.third_interest, 0, 0, 1)
        ) AS t(faculty, first_count, second_count, third_count)
        WHERE t.faculty IS NOT NULL
        GROUP BY t.faculty
        ORDER BY (SUM(t.first_count) + SUM(t.second_count) + SUM(t.third_count)) DESC;
    `
	err := r.DB.Raw(query).Scan(&results).Error
	return results, err
}

func (r *DashBoardRepository) GetSourceCount() ([]domain.SourceCount, error) {
	var results []domain.SourceCount
	query :=
		`SELECT source, COUNT(*) as count FROM (
            SELECT unnest(selected_sources) as source FROM users WHERE selected_sources IS NOT NULL
        ) AS sources
        GROUP BY source
        ORDER BY count DESC;`
	err := r.DB.Raw(query).Scan(&results).Error
	return results, err
}

func (r *DashBoardRepository) GetAgeGroupCount() ([]domain.AgeCount, error) {
	var results []domain.AgeCount

	err := r.DB.Model(&domain.User{}).
		Select("EXTRACT(YEAR FROM AGE(birth_date)) AS age, COUNT(*) AS count").
		Where("birth_date IS NOT NULL").
		Group("age").
		Order("age ASC").
		Scan(&results).Error

	return results, err
}

func (r *DashBoardRepository) GetFacultyToday() (domain.FacultyRegisterCount, error) {
	var result domain.FacultyRegisterCount

	// 1. เขียน Query หาคณะที่ลงทะเบียนมากที่สุดในวันนี้
	err := r.DB.Model(&domain.StudentTransaction{}).
		Select("faculty, COUNT(*) as count").
		Where("DATE(registered_at) = CURRENT_DATE").
		Group("faculty").
		Order("count DESC").
		Scan(&result).Error

	return result, err
}

func (r *DashBoardRepository) GetStatusStudent() ([]domain.StatusCount, error) {
	var results []domain.StatusCount
	query :=
		`SELECT status, COUNT(*) as count FROM users
		WHERE status IS NOT NULL
		GROUP BY status
		ORDER BY count DESC;`
	err := r.DB.Raw(query).Scan(&results).Error
	return results, err
}
