package repository

import (
	"github.com/isd-sgcu/oph-67-backend/domain"
	"gorm.io/gorm"
)

type StudentEvaluationRepository struct {
	DB *gorm.DB
}

func NewStudentEvaluationRepository(db *gorm.DB) *StudentEvaluationRepository {
	return &StudentEvaluationRepository{DB: db}
}

func (r *StudentEvaluationRepository) CreateStudentEvaluation(evaluation *domain.StudentEvaluation) error {
	return r.DB.Create(evaluation).Error
}

func (r *StudentEvaluationRepository) GetStudentEvaluationByStudentId(studentId string) (*domain.StudentEvaluation, error) {
	var evaluation domain.StudentEvaluation
	err := r.DB.Where("student_id = ?", studentId).First(&evaluation).Error
	if err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (r *StudentEvaluationRepository) UpdateStudentEvaluation(evaluation *domain.StudentEvaluation) error {
	return r.DB.Save(evaluation).Error
}

func (r *StudentEvaluationRepository) DeleteStudentEvaluation(studentId string) error {
	return r.DB.Where("student_id = ?", studentId).Delete(&domain.StudentEvaluation{}).Error
}

func (r *StudentEvaluationRepository) GetAllStudentEvaluations() ([]domain.StudentEvaluation, error) {
	var evaluations []domain.StudentEvaluation
	err := r.DB.Find(&evaluations).Error
	if err != nil {
		return nil, err
	}
	return evaluations, nil
}
func (r *StudentEvaluationRepository) GetStudentEvaluationCount() (int64, error) {
	var count int64
	err := r.DB.Model(&domain.StudentEvaluation{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StudentEvaluationRepository) GetStudentEvaluationById(id string) (*domain.StudentEvaluation, error) {
	var evaluation domain.StudentEvaluation
	err := r.DB.Where("id = ?", id).First(&evaluation).Error
	if err != nil {
		return nil, err
	}
	return &evaluation, nil
}
