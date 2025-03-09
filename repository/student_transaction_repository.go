package repository

import (
	"github.com/isd-sgcu/oph-67-backend/domain"
	"gorm.io/gorm"
)

type StudentTransactionRepository struct {
	DB *gorm.DB
}

func NewStudentTransactionRepository(db *gorm.DB) *StudentTransactionRepository {
	return &StudentTransactionRepository{DB: db}
}

// Create a new student transaction
func (r *StudentTransactionRepository) Create(transaction *domain.StudentTransaction) error {
	return r.DB.Create(transaction).Error
}

// Get all student transactions
func (r *StudentTransactionRepository) GetAll() ([]domain.StudentTransaction, error) {
	var transactions []domain.StudentTransaction
	err := r.DB.Find(&transactions).Error
	return transactions, err
}

// Get a student transaction by ID
func (r *StudentTransactionRepository) GetById(id string) (domain.StudentTransaction, error) {
	var transaction domain.StudentTransaction
	err := r.DB.Where("id = ?", id).First(&transaction).Error
	return transaction, err
}

// Get student transactions by student ID
func (r *StudentTransactionRepository) GetByStudentId(studentId string) ([]domain.StudentTransaction, error) {
	var transactions []domain.StudentTransaction
	err := r.DB.Where("student_registration_id = ?", studentId).Find(&transactions).Error
	return transactions, err
}

// Update a student transaction
func (r *StudentTransactionRepository) Update(id string, transaction *domain.StudentTransaction) error {
	err := r.DB.Model(&domain.StudentTransaction{}).Where("id = ?", id).Updates(transaction).Error
	return err
}

// Delete a student transaction by ID
func (r *StudentTransactionRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&domain.StudentTransaction{}).Error
	return err
}

func (r *StudentTransactionRepository) GetByStudentIdAndFaculty(studentId string, faculty string) ([]domain.StudentTransaction, error) {
	var transactions []domain.StudentTransaction
	err := r.DB.Where("student_registration_id = ? AND faculty = ?", studentId, faculty).Find(&transactions).Error
	return transactions, err
}
