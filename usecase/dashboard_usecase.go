package usecase

import (
	"github.com/isd-sgcu/oph-67-backend/domain"
)

type DashboardUseCase struct {
	DashboardRepo DashBoardRepositoryInterface
}

type DashBoardRepositoryInterface interface {
	GetFacultyCount() ([]domain.FacultyPercent, error)
	GetSourceCount() ([]domain.SourceCount, error)
}

func NewDashBoardUseCase(dashboardRepo DashBoardRepositoryInterface) *DashboardUseCase {
	return &DashboardUseCase{DashboardRepo: dashboardRepo}
}

func (d *DashboardUseCase) GetFacultyCount() ([]domain.FacultyPercent, error) {
	return d.DashboardRepo.GetFacultyCount()
}

func (d *DashboardUseCase) GetSourceCount() ([]domain.SourceCount, error) {
	return d.DashboardRepo.GetSourceCount()
}
