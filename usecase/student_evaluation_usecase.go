package usecase

import (
	"strconv"

	"github.com/isd-sgcu/oph-67-backend/domain"
)

type StudentEvaluationUsecase struct {
	StudentEvaluationRepo StudentEvaluationRepositoryInterface
}

type StudentEvaluationRepositoryInterface interface {
	CreateStudentEvaluation(evaluation *domain.StudentEvaluation) error
	GetStudentEvaluationByStudentId(studentId string) (*domain.StudentEvaluation, error)
	UpdateStudentEvaluation(evaluation *domain.StudentEvaluation) error
	DeleteStudentEvaluation(studentId string) error
	GetAllStudentEvaluations() ([]domain.StudentEvaluation, error)
	GetStudentEvaluationCount() (int64, error)
	GetStudentEvaluationById(id string) (*domain.StudentEvaluation, error)
}

func NewStudentEvaluationUsecase(studentEvaluationRepo StudentEvaluationRepositoryInterface) *StudentEvaluationUsecase {
	return &StudentEvaluationUsecase{StudentEvaluationRepo: studentEvaluationRepo}
}

func (u *StudentEvaluationUsecase) CreateStudentEvaluation(evaluation *domain.StudentEvaluation) error {
	isExist, _ := u.StudentEvaluationRepo.GetStudentEvaluationByStudentId(evaluation.StudentId)

	if isExist != nil {
		return domain.ErrStudentEvaluationAlreadyExists
	}
	return u.StudentEvaluationRepo.CreateStudentEvaluation(evaluation)
}

func (u *StudentEvaluationUsecase) GetStudentEvaluationByStudentId(studentId string) (*domain.StudentEvaluation, error) {
	return u.StudentEvaluationRepo.GetStudentEvaluationByStudentId(studentId)
}

func (u *StudentEvaluationUsecase) UpdateStudentEvaluation(evaluation *domain.StudentEvaluation) error {
	return u.StudentEvaluationRepo.UpdateStudentEvaluation(evaluation)
}

func (u *StudentEvaluationUsecase) DeleteStudentEvaluation(studentId string) error {
	return u.StudentEvaluationRepo.DeleteStudentEvaluation(studentId)
}

func (u *StudentEvaluationUsecase) GetAllStudentEvaluations() ([]domain.StudentEvaluation, error) {
	return u.StudentEvaluationRepo.GetAllStudentEvaluations()
}

func (u *StudentEvaluationUsecase) GetStudentEvaluationCount() (int64, error) {
	return u.StudentEvaluationRepo.GetStudentEvaluationCount()
}

func (u *StudentEvaluationUsecase) GetStudentEvaluationById(id string) (*domain.StudentEvaluation, error) {
	return u.StudentEvaluationRepo.GetStudentEvaluationById(id)
}

func (u *StudentEvaluationUsecase) GetStudentEvaluationByStudentIdAndId(studentId string, id string) (*domain.StudentEvaluation, error) {
	evaluation, err := u.StudentEvaluationRepo.GetStudentEvaluationByStudentId(studentId)
	if err != nil {
		return nil, err
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if evaluation.ID != idInt {
		return nil, domain.ErrUserNotFound
	}
	return evaluation, nil
}
