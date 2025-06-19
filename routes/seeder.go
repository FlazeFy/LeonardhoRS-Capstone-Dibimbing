package routes

import (
	"pelita/repository"
	"pelita/seeder"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, roomRepo repository.RoomRepository, adminRepo repository.AdminRepository, technicianRepo repository.TechnicianRepository, userRepo repository.UserRepository, assetRepo repository.AssetRepository, assetPlacement repository.AssetPlacementRepository, assetMaintenance repository.AssetMaintenanceRepository) {
	seeder.SeedRooms(roomRepo, 20)
	seeder.SeedAdmins(adminRepo, 5)
	seeder.SeedTechnicians(technicianRepo, adminRepo, 40)
	seeder.SeedUsers(userRepo, 80)
	seeder.SeedAssets(assetRepo, adminRepo, 200)
	seeder.SeedAssetPlacements(assetPlacement, adminRepo, roomRepo, assetRepo, technicianRepo, 350)
	seeder.SeedAssetMaintenances(assetMaintenance, adminRepo, technicianRepo, assetPlacement, 500)
}
