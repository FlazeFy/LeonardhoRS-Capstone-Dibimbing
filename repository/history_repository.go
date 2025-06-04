package repository

import (
	"errors"
	"fmt"
	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HistoryRepository interface {
	FindAll() ([]entity.AllHistory, error)
	FindMy(id uuid.UUID, typeUser string) ([]entity.History, error)
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

func (r *historyRepository) FindMy(id uuid.UUID, typeUser string) ([]entity.History, error) {
	// Models
	var history []entity.History

	// Query
	var targetCol string
	if typeUser == "admin" {
		targetCol = "admin_id"
	} else if typeUser == "technician" {
		targetCol = "technician_id"
	} else if typeUser == "guest" {
		targetCol = "user_id"
	}

	err := r.db.Where(fmt.Sprintf("%s = ?", targetCol), id).
		Where("type_user = ?", typeUser).
		Order("created_at DESC").
		Find(&history).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return history, err
}
