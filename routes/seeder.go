package routes

import (
	"pelita/repository"
	"pelita/seeder"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, roomRepo repository.RoomRepository, adminRepo repository.AdminRepository, technicianRepo repository.TechnicianRepository, userRepo repository.UserRepository, assetRepo repository.AssetRepository, assetPlacement repository.AssetPlacementRepository) {
	seeder.SeedRooms(roomRepo, 20)
	seeder.SeedAdmins(adminRepo, 5)
	seeder.SeedTechnicians(technicianRepo, adminRepo, 40)
	seeder.SeedUsers(userRepo, 80)
	seeder.SeedAssets(assetRepo, adminRepo, 250)
	seeder.SeedAssetPlacements(assetPlacement, adminRepo, roomRepo, assetRepo, technicianRepo, 500)
}
