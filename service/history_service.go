package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

type HistoryService interface {
	GetAllHistory(pagination utils.Pagination) ([]entity.AllHistory, int64, error)
	GetMyHistory(pagination utils.Pagination, id uuid.UUID, typeUser string) ([]entity.History, int64, error)
}

type historyService struct {
	historyRepo repository.HistoryRepository
}

func NewHistoryService(historyRepo repository.HistoryRepository) HistoryService {
	return &historyService{
		historyRepo: historyRepo,
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
