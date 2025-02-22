package repository

import (
	"github.com/isd-sgcu/oph-67-backend/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetById(id string) (domain.User, error) {
	var user domain.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *UserRepository) GetByName(name string) ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Where("name ILIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByPhone(phone string) (domain.User, error) {
	var user domain.User
	err := r.DB.Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (r *UserRepository) Update(id string, user *domain.User) error {
	err := r.DB.Model(&domain.User{}).Where("id = ?", id).Updates(user).Error
	return err
}

func (r *UserRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&domain.User{}).Error
	return err
}

func (r *UserRepository) IsUIDExists(uid string) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.User{}).Where("uid = ?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
