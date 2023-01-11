package repositories

import (
	"server/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUsers() ([]models.User, error)
	GetUser(ID int) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User, ID int) (models.User, error)
	DeleteUser(user models.User, ID int) (models.User, error)
}

func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error // ORM
	return users, err
}

func (r *repository) GetUser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error // ORM

	return user, err
}
func (r *repository) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error // ORM
	return user, err
}

func (r *repository) UpdateUser(user models.User, ID int) (models.User, error) {
	err := r.db.Save(&user).Error // ORM
	return user, err
}

func (r *repository) DeleteUser(user models.User, ID int) (models.User, error) {
	err := r.db.Delete(&user).Error //ORM
	return user, err
}
