package scheduler

import (
	"fmt"
	"log"
	"os"
	"pelita/service"
	"pelita/utils"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type AssetMaintenanceScheduler struct {
	AssetMaintenanceService service.AssetMaintenanceService
	AssetFindingService     service.AssetFindingService
	AdminService            service.AdminService
}

func NewAssetMaintenanceScheduler(
	assetMaintenanceService service.AssetMaintenanceService,
	assetFindingService service.AssetFindingService,
	adminService service.AdminService,
) *AssetMaintenanceScheduler {
	return &AssetMaintenanceScheduler{
		AssetMaintenanceService: assetMaintenanceService,
		AssetFindingService:     assetFindingService,
		AdminService:            adminService,
	}
}

func (s *AssetMaintenanceScheduler) ReminderSchedulerTodayMaintenance() {
	// Service : Get Today Schedule Maintenance
	scheduleMap, err := s.AssetMaintenanceService.GetTodayValidSchedules()
	if err != nil {
		log.Println("Failed to fetch today's maintenance schedules:", err)
		return
	}

	// Service : Get All Admin Contact
	adminContacts, err := s.AdminService.GetAllContact()
	if err != nil {
		log.Println("Failed to fetch admin contacts:", err)
		return
	}

	if len(scheduleMap) == 0 {
		log.Println("No maintenance schedule today.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Println("Failed to connect to Telegram bot:", err)
		return
	}

	// Admin Message
	fullMessage := "🛠 *Today's All Maintenance Schedule:*\n\n"
	index := 1
	for _, schedules := range scheduleMap {
		for _, s := range schedules {
			fullMessage += fmt.Sprintf("%d. %s (%s)\n⏰ %s - %s\n👨 %s\n📧 %s\n📝 %s\n\n",
				index,
				s.AssetName,
				s.AssetCategory,
				s.MaintenanceHourStart.Format("15:04"),
				s.MaintenanceHourEnd.Format("15:04"),
				s.Username,
				s.Email,
				utils.NullSafeString(s.MaintenanceNotes),
			)
			index++
		}
	}

	// Send Admin Message
	for _, contact := range adminContacts {
		if contact.TelegramUserId == nil || !contact.TelegramIsValid {
			continue
		}

		telegramID, err := strconv.ParseInt(*contact.TelegramUserId, 10, 64)
		if err != nil {
			log.Printf("Invalid Telegram ID for admin %s: %v\n", contact.Username, err)
			continue
		}

		msg := tgbotapi.NewMessage(telegramID, fullMessage)
		msg.ParseMode = "Markdown"

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send message to admin %s: %v\n", contact.Username, err)
		} else {
			log.Printf("Full schedule sent to admin %s (%s)\n", contact.Username, *contact.TelegramUserId)
		}
	}

	// Send Technician Message
	for telegramUserID, schedules := range scheduleMap {
		if telegramUserID == "" || len(schedules) == 0 {
			continue
		}

		telegramID, err := strconv.ParseInt(telegramUserID, 10, 64)
		if err != nil {
			log.Printf("Invalid technician Telegram ID: %v\n", err)
			continue
		}

		personalMessage := "🛠 *Your Maintenance Schedule for Today:*\n\n"
		for i, s := range schedules {
			personalMessage += fmt.Sprintf("%d. %s (%s)\n⏰ %s - %s\n📝 %s\n\n",
				i+1,
				s.AssetName,
				s.AssetCategory,
				s.MaintenanceHourStart.Format("15:04"),
				s.MaintenanceHourEnd.Format("15:04"),
				utils.NullSafeString(s.MaintenanceNotes),
			)
		}

		msg := tgbotapi.NewMessage(telegramID, personalMessage)
		msg.ParseMode = "Markdown"

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send schedule to technician %s: %v\n", telegramUserID, err)
		} else {
			log.Printf("Personal schedule sent to technician (%s)\n", telegramUserID)
		}
	}
}

func (s *AssetMaintenanceScheduler) AuditSchedulerAssetFindingReport() {
	// Service : Get All Asset Finding Report
	findingMap, err := s.AssetFindingService.GetAllAssetFindingReport()
	if err != nil {
		log.Println("Failed to fetch today's maintenance schedules:", err)
		return
	}

	// Service : Get All Admin Contact
	adminContacts, err := s.AdminService.GetAllContact()
	if err != nil {
		log.Println("Failed to fetch admin contacts:", err)
		return
	}

	if len(findingMap) == 0 {
		log.Println("No finding found.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Println("Failed to connect to Telegram bot:", err)
		return
	}

	// Asset Finding Docs
	datetime := time.Now()
	filename := fmt.Sprintf("audit_asset_finding_%s.pdf", datetime)
	err = utils.GeneratePDFAssetFindingReport(findingMap, filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send Admin Message
	for _, contact := range adminContacts {
		if contact.TelegramUserId == nil || !contact.TelegramIsValid {
			continue
		}

		telegramID, err := strconv.ParseInt(*contact.TelegramUserId, 10, 64)
		if err != nil {
			log.Printf("Invalid Telegram ID for admin %s: %v\n", contact.Username, err)
			continue
		}

		msg := tgbotapi.NewDocumentUpload(telegramID, filename)
		msg.ParseMode = "Markdown"
		msg.Caption = fmt.Sprintf("[ADMIN] Hello %s, We're here to report asset finding that we've got so far. Here are the docs:", contact.Username)

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send message to admin %s: %v\n", contact.Username, err)
		} else {
			log.Printf("Full schedule sent to admin %s (%s)\n", contact.Username, *contact.TelegramUserId)
		}
	}
}
