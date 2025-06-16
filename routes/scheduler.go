package routes

import (
	"pelita/scheduler"
	"pelita/service"
	"time"

	"github.com/robfig/cron"
)

func SetUpScheduler(assetMaintenanceService service.AssetMaintenanceService, assetFindingService service.AssetFindingService, adminService service.AdminService) {
	// Initialize Scheduler
	maintenanceScheduler := scheduler.NewAssetMaintenanceScheduler(assetMaintenanceService, assetFindingService, adminService)

	// Init Scheduler
	c := cron.New()
	Scheduler(c, maintenanceScheduler)
	c.Start()
	defer c.Stop()
}

func Scheduler(c *cron.Cron, maintenanceScheduler *scheduler.AssetMaintenanceScheduler) {
	// Production
	// Every day at 00:10 AM)
	c.AddFunc("10 0 * * *", maintenanceScheduler.ReminderSchedulerTodayMaintenance)
	// Every day at 00:20 AM)
	c.AddFunc("20 0 * * *", maintenanceScheduler.AuditSchedulerAssetFindingReport)

	// Development (after 5 sec)
	go func() {
		time.Sleep(5 * time.Second)
		maintenanceScheduler.ReminderSchedulerTodayMaintenance()
		maintenanceScheduler.AuditSchedulerAssetFindingReport()
	}()
}
