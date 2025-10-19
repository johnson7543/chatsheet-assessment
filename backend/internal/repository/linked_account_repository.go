package repository

import (
	"github.com/johnson7543/chatsheet-assessment/internal/models"
	"gorm.io/gorm"
)

// LinkedAccountRepository handles linked account data operations
type LinkedAccountRepository struct {
	db *gorm.DB
}

// NewLinkedAccountRepository creates a new linked account repository
func NewLinkedAccountRepository(db *gorm.DB) *LinkedAccountRepository {
	return &LinkedAccountRepository{db: db}
}

// Create creates a new linked account
func (r *LinkedAccountRepository) Create(account *models.LinkedAccount) error {
	return r.db.Create(account).Error
}

// FindByUserID finds all linked accounts for a user
func (r *LinkedAccountRepository) FindByUserID(userID uint) ([]models.LinkedAccount, error) {
	var accounts []models.LinkedAccount
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&accounts).Error
	return accounts, err
}

// FindByID finds a linked account by ID
func (r *LinkedAccountRepository) FindByID(id uint) (*models.LinkedAccount, error) {
	var account models.LinkedAccount
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindByUserIDAndID finds a linked account by ID and user ID (for authorization)
func (r *LinkedAccountRepository) FindByUserIDAndID(userID, id uint) (*models.LinkedAccount, error) {
	var account models.LinkedAccount
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Delete soft deletes a linked account
func (r *LinkedAccountRepository) Delete(account *models.LinkedAccount) error {
	return r.db.Delete(account).Error
}

// CountByUserID counts linked accounts for a user
func (r *LinkedAccountRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.LinkedAccount{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
