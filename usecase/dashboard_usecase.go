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
	GetFacultyToday() ([]domain.FacultyRegisterCount, error)
	GetStatusStudent() ([]domain.StatusCount, error)
	GetAllStudents() ([]domain.StudentProfile, error)
	GetStudentsByFacultyInterest(faculty string) ([]domain.StudentProfile, error)
	GetAttendedCount() ([]domain.AttendedCount, error)
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

func (d *DashboardUseCase) GetFacultyTodayCount() ([]domain.FacultyRegisterCount, error) {
	return d.DashboardRepo.GetFacultyToday()
}

func (d *DashboardUseCase) GetStatusStudent() ([]domain.StatusCount, error) {
	return d.DashboardRepo.GetStatusStudent()
}

func (d *DashboardUseCase) GetAllStudent() ([]domain.StudentProfile, error) {
	return d.DashboardRepo.GetAllStudents()
}

func (d *DashboardUseCase) GetStudentsByFacultyInterest(faculty string) ([]domain.StudentProfile, error) {
	return d.DashboardRepo.GetStudentsByFacultyInterest(faculty)
}

func (d *DashboardUseCase) GetAttendedCount() ([]domain.AttendedCount, error) {
	return d.DashboardRepo.GetAttendedCount()
}
