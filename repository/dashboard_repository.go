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
            SUM(t.third_count) AS third_count
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
