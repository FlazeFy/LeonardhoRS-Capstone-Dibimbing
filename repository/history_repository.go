package repository

import (
	"errors"
	"pelita/entity"

	"gorm.io/gorm"
)

type HistoryRepository interface {
	FindAll() ([]entity.AllHistory, error)
}

type historyRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) FindAll() ([]entity.AllHistory, error) {
	// Models
	var history []entity.AllHistory

	// Query
	err := r.db.Table("histories").
		Preload("User").
		Preload("Technician").
		Preload("Admin").
		Order("created_at DESC").
		Find(&history).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return history, err
}
