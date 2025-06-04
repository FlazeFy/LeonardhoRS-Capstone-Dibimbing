package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type HistoryService interface {
	GetAllHistory() ([]entity.AllHistory, error)
	GetMyHistory(id uuid.UUID, typeUser string) ([]entity.History, error)
}

type historyService struct {
	historyRepo repository.HistoryRepository
}

func NewHistoryService(historyRepo repository.HistoryRepository) HistoryService {
	return &historyService{
		historyRepo: historyRepo,
	}
}

func (s *historyService) GetAllHistory() ([]entity.AllHistory, error) {
	// Repo : Get All History
	history, err := s.historyRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if history == nil {
		return nil, errors.New("history not found")
	}

	return history, nil
}

func (s *historyService) GetMyHistory(id uuid.UUID, typeUser string) ([]entity.History, error) {
	// Repo : Get My History
	history, err := s.historyRepo.FindMy(id, typeUser)
	if err != nil {
		return nil, err
	}
	if history == nil {
		return nil, errors.New("history not found")
	}

	return history, nil
}
