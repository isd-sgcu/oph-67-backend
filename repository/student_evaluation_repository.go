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

// Create a new student evaluation
func (r *StudentEvaluationRepository) Create(evaluation *domain.StudentEvaluation) error {
	return r.DB.Create(evaluation).Error
}

// Get all student evaluations
func (r *StudentEvaluationRepository) GetAll() ([]domain.StudentEvaluation, error) {
	var evaluations []domain.StudentEvaluation
	err := r.DB.Find(&evaluations).Error
	return evaluations, err
}

// Get a student evaluation by ID
func (r *StudentEvaluationRepository) GetById(id string) (domain.StudentEvaluation, error) {
	var evaluation domain.StudentEvaluation
	err := r.DB.Where("id = ?", id).First(&evaluation).Error
	return evaluation, err
}

// Get student evaluations by student ID
func (r *StudentEvaluationRepository) GetByStudentId(studentId string) ([]domain.StudentEvaluation, error) {
	var evaluations []domain.StudentEvaluation
	err := r.DB.Where("student_registration_id = ?", studentId).Find(&evaluations).Error
	return evaluations, err
}

// Update a student evaluation
func (r *StudentEvaluationRepository) Update(id string, evaluation *domain.StudentEvaluation) error {
	err := r.DB.Model(&domain.StudentEvaluation{}).Where("id = ?", id).Updates(evaluation).Error
	return err
}

// Delete a student evaluation by ID
func (r *StudentEvaluationRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&domain.StudentEvaluation{}).Error
	return err
}