package repositories

import (
	"server/models"

	"gorm.io/gorm"
)

type DonationRepository interface {
	FindDonations() ([]models.Donation, error)
	GetDonation(ID int) (models.Donation, error)
	CreateDonation(donation models.Donation) (models.Donation, error)
	UpdateDonation(donation models.Donation, ID int) (models.Donation, error)
	DeleteDonation(donation models.Donation, ID int) (models.Donation, error)
}

func RepositoryDonation(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindDonations() ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Find(&donations).Error
	return donations, err
}

func (r *repository) GetDonation(ID int) (models.Donation, error) {
	var donation models.Donation
	err := r.db.First(&donation, ID).Error
	return donation, err
}

func (r *repository) CreateDonation(donation models.Donation) (models.Donation, error) {
	err := r.db.Preload("Funding").Create(&donation).Error
	return donation, err
}

func (r *repository) UpdateDonation(donation models.Donation, ID int) (models.Donation, error) {
	err := r.db.Preload("Funding").Save(&donation).Error
	return donation, err
}

func (r *repository) DeleteDonation(donation models.Donation, ID int) (models.Donation, error) {
	err := r.db.Delete(&donation).Error
	return donation, err
}
