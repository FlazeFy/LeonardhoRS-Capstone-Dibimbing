package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

type AssetMaintenanceService interface {
	GetAllAssetMaintenance(pagination utils.Pagination) ([]entity.AssetMaintenance, int64, error)
	GetAllAssetMaintenanceSchedule() ([]entity.AssetMaintenanceSchedule, error)
	GetMostContext(targetCol string) ([]entity.StatsContextTotal, error)
	Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error
	UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error
	DeleteById(id uuid.UUID) error

	// Scheduler Service
	GetTodayValidSchedules() (map[string][]entity.AssetMaintenanceSchedule, error)
}

type assetMaintenanceService struct {
	assetMaintenanceRepo repository.AssetMaintenanceRepository
	technicianRepo       repository.TechnicianRepository
	assetRepo            repository.AssetRepository
	statsRepo            repository.StatsRepository
}

func NewAssetMaintenanceService(assetMaintenanceRepo repository.AssetMaintenanceRepository, technicianRepo repository.TechnicianRepository, assetRepo repository.AssetRepository, statsRepo repository.StatsRepository) AssetMaintenanceService {
	return &assetMaintenanceService{
		assetMaintenanceRepo: assetMaintenanceRepo,
		technicianRepo:       technicianRepo,
		assetRepo:            assetRepo,
		statsRepo:            statsRepo,
	}
}

func (s *assetMaintenanceService) GetAllAssetMaintenance(pagination utils.Pagination) ([]entity.AssetMaintenance, int64, error) {
	// Repo : Get All Asset Maintenance
	assetMaintenance, total, err := s.assetMaintenanceRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if assetMaintenance == nil {
		return nil, 0, errors.New("asset maintenance not found")
	}

	return assetMaintenance, total, nil
}

func (s *assetMaintenanceService) GetAllAssetMaintenanceSchedule() ([]entity.AssetMaintenanceSchedule, error) {
	// Repo : Get All Asset Maintenance Schedule
	assetMaintenance, err := s.assetMaintenanceRepo.FindAllSchedule()
	if err != nil {
		return nil, err
	}
	if assetMaintenance == nil {
		return nil, errors.New("asset maintenance schedule not found")
	}

	return assetMaintenance, nil
}

func (s *assetMaintenanceService) Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error {
	// Validator
	if assetMaintenance.AssetPlacementId == uuid.Nil {
		return errors.New("asset placement id is required")
	}
	if assetMaintenance.MaintenanceBy == uuid.Nil {
		return errors.New("asset maintenance by is required")
	}
	if assetMaintenance.MaintenanceDay == "" {
		return errors.New("asset maintenance day is required")
	}

	// Repo : Get Asset Maintenance By Asset Placement Id, Maintenance By, and Maintenance Day
	is_exist, err := s.assetMaintenanceRepo.FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetMaintenance.AssetPlacementId, assetMaintenance.MaintenanceBy, assetMaintenance.MaintenanceDay, assetMaintenance.MaintenanceHourStart, assetMaintenance.MaintenanceHourEnd)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset is already assigned to maintenance by a technician")
	}

	// Repo : Create Asset Maintenance
	if err := s.assetMaintenanceRepo.Create(assetMaintenance, adminId); err != nil {
		return err
	}

	// Repo : Find Technician By Id
	technician, err := s.technicianRepo.FindById(assetMaintenance.MaintenanceBy)
	if err != nil {
		return err
	}

	// Send Telegram
	if technician != nil && technician.TelegramIsValid && technician.TelegramUserId != nil {
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to connect to Telegram bot:", err.Error()))
		}

		telegramID, err := strconv.ParseInt(*technician.TelegramUserId, 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid technician Telegram ID: %v\n", err))
		}

		// Repo : Find Asset By Id
		asset, err := s.assetRepo.FindByAssetPlacementId(assetMaintenance.AssetPlacementId)
		if err != nil {
			return err
		}

		// Build Message
		personalMessage := fmt.Sprintf("🛠 *You Have A New Asset To Maintenance:*\n\nAsset Name : %s\nCategory : %s\nDay / Hour : Every %s at %s - %s\n",
			asset.AssetName,
			asset.AssetCategory,
			assetMaintenance.MaintenanceDay,
			assetMaintenance.MaintenanceHourStart.Format("15:04"),
			assetMaintenance.MaintenanceHourEnd.Format("15:04"))

		msg := tgbotapi.NewMessage(telegramID, personalMessage)
		msg.ParseMode = "Markdown"

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send message to technician %s: %v\n", *technician.TelegramUserId, err)
		} else {
			log.Printf("Personal message sent to technician (%s)\n", *technician.TelegramUserId)
		}
	}

	return nil
}

func (s *assetMaintenanceService) UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error {
	// Validator
	if assetMaintenance.AssetPlacementId == uuid.Nil {
		return errors.New("asset placement id is required")
	}
	if assetMaintenance.MaintenanceBy == uuid.Nil {
		return errors.New("asset maintenance by is required")
	}
	if assetMaintenance.MaintenanceDay == "" {
		return errors.New("asset maintenance day is required")
	}

	// Repo : Get Asset Maintenance By Asset Placement Id, Maintenance By, Maintenance Day, and Id
	is_exist, err := s.assetMaintenanceRepo.FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetMaintenance.AssetPlacementId, assetMaintenance.MaintenanceBy, assetMaintenance.MaintenanceDay, assetMaintenance.MaintenanceHourStart, assetMaintenance.MaintenanceHourEnd, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset is already assigned to maintenance by a technician")
	}

	// Repo : Update Asset Maintenance By Id
	if err := s.assetMaintenanceRepo.UpdateById(assetMaintenance, id); err != nil {
		return err
	}

	return nil
}

func (s *assetMaintenanceService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Asset Maintenance By Id
	err := s.assetMaintenanceRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *assetMaintenanceService) GetMostContext(targetCol string) ([]entity.StatsContextTotal, error) {
	// Repo : Get My History
	asset, err := s.statsRepo.FindMostUsedContext("asset_maintenances", targetCol)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset maintenance not found")
	}

	return asset, nil
}

// Scheduler Service
func (s *assetMaintenanceService) GetTodayValidSchedules() (map[string][]entity.AssetMaintenanceSchedule, error) {
	allSchedules, err := s.assetMaintenanceRepo.FindAllSchedule()
	if err != nil {
		return nil, err
	}

	today := time.Now().Weekday().String()[:3]

	// Group by telegram_user_id
	result := make(map[string][]entity.AssetMaintenanceSchedule)
	for _, schedule := range allSchedules {
		if schedule.MaintenanceDay == today &&
			schedule.TelegramUserId != nil &&
			schedule.TelegramIsValid {
			result[*schedule.TelegramUserId] = append(result[*schedule.TelegramUserId], schedule)
		}
	}

	return result, nil
}
