package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

// History Interface
type HistoryService interface {
	GetAllHistory(pagination utils.Pagination) ([]entity.AllHistory, int64, error)
	GetMyHistory(pagination utils.Pagination, id uuid.UUID, typeUser string) ([]entity.History, int64, error)
	GetMostContext(targetCol string) ([]entity.StatsContextTotal, error)
}

// History Struct
type historyService struct {
	historyRepo repository.HistoryRepository
	statsRepo   repository.StatsRepository
}

// History Constructor
func NewHistoryService(historyRepo repository.HistoryRepository, statsRepo repository.StatsRepository) HistoryService {
	return &historyService{
		historyRepo: historyRepo,
		statsRepo:   statsRepo,
	}
}

func (s *historyService) GetAllHistory(pagination utils.Pagination) ([]entity.AllHistory, int64, error) {
	// Repo : Get All History
	history, total, err := s.historyRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if history == nil {
		return nil, 0, errors.New("history not found")
	}

	return history, total, nil
}

func (s *historyService) GetMyHistory(pagination utils.Pagination, id uuid.UUID, typeUser string) ([]entity.History, int64, error) {
	// Repo : Get My History
	history, total, err := s.historyRepo.FindMy(pagination, id, typeUser)
	if err != nil {
		return nil, 0, err
	}
	if history == nil {
		return nil, 0, errors.New("history not found")
	}

	return history, total, nil
}

func (s *historyService) GetMostContext(targetCol string) ([]entity.StatsContextTotal, error) {
	// Repo : Get My History
	history, err := s.statsRepo.FindMostUsedContext("histories", targetCol)
	if err != nil {
		return nil, err
	}
	if history == nil {
		return nil, errors.New("history not found")
	}

	return history, nil
}
