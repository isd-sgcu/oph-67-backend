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
	GetAgeGroupCount() ([]domain.AgeCount, error)
	GetFacultyTodayCount() ([]domain.FacultyPercent, error)
	GetStatusStudent()	([]domain.StatusCount, error)
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

func (d *DashboardUseCase) GetAgeGroupCount() ([]domain.AgeCount, error) {
	return d.DashboardRepo.GetAgeGroupCount()
}

func (d *DashboardUseCase) GetFacultyTodayCount() ([]domain.FacultyPercent, error) {
	return d.DashboardRepo.GetFacultyCount()
}

func (d *DashboardUseCase) GetStatusStudent() ([]domain.SourceCount, error) {
	return d.DashboardRepo.GetSourceCount()
}
