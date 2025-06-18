package routes

import (
	"pelita/repository"
	"pelita/seeder"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, roomRepo repository.RoomRepository, adminRepo repository.AdminRepository) {
	seeder.SeedRooms(roomRepo, 10)
	seeder.SeedAdmins(adminRepo, 5)
}
