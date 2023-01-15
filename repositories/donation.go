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
	GetDonationByUser(ID int) ([]models.Donation, error)
	GetDonationByFunding(ID int) ([]models.Donation, error)
	GetDonationPending(ID int) ([]models.Donation, error)
	UpdateStatus(status string, ID int) error
}

func RepositoryDonation(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindDonations() ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Preload("User").Preload("Funding").Find(&donations).Error
	return donations, err
}

func (r *repository) GetDonation(ID int) (models.Donation, error) {
	var donation models.Donation
	err := r.db.Preload("Funding").Preload("User").First(&donation, ID).Error
	return donation, err
}

func (r *repository) CreateDonation(donation models.Donation) (models.Donation, error) {
	err := r.db.Create(&donation).Error
	return donation, err
}

func (r *repository) UpdateDonation(donation models.Donation, ID int) (models.Donation, error) {
	err := r.db.Save(&donation).Error
	return donation, err
}

func (r *repository) DeleteDonation(donation models.Donation, ID int) (models.Donation, error) {
	err := r.db.Delete(&donation).Error
	return donation, err
}

func (r *repository) GetDonationByUser(ID int) ([]models.Donation, error) {
	var donation []models.Donation
	err := r.db.Preload("User").Preload("Funding").Where("user_id = ?", ID).Find(&donation).Error
	return donation, err
}
func (r *repository) GetDonationByFunding(ID int) ([]models.Donation, error) {
	var donation []models.Donation
	err := r.db.Preload("User").Preload("Funding").Where("funding_id = ?", ID).Where("status = 'success'").Find(&donation).Error
	return donation, err
}
func (r *repository) GetDonationPending(ID int) ([]models.Donation, error) {
	var donation []models.Donation
	err := r.db.Preload("User").Preload("Funding").Where("funding_id = ?", ID).Where("status <> 'success'").Find(&donation).Error
	return donation, err
}
func (r *repository) UpdateStatus(status string, ID int) error {
	var donation models.Donation
	r.db.Preload("Funding").First(&donation, ID)

	donation.Status = status

	err := r.db.Save(&donation).Error

	return err
}
