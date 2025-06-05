package repository

import (
	"errors"
	"fmt"
	"pelita/entity"

	"gorm.io/gorm"
)

type StatsRepository interface {
	FindMostUsedContext(tableName, targetCol string) ([]entity.StatsContextTotal, error)
}

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) FindMostUsedContext(tableName, targetCol string) ([]entity.StatsContextTotal, error) {

	// Models
	var stats []entity.StatsContextTotal

	// Query
	err := r.db.Table(tableName).
		Select(fmt.Sprintf("COUNT(%s) as total, %s as context", targetCol, targetCol)).
		Group(targetCol).
		Order("total DESC").
		Limit(7).
		Find(&stats).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return stats, nil
}
