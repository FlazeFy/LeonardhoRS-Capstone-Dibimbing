package routes

import (
	"pelita/repository"
	"pelita/seeder"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, roomRepo repository.RoomRepository, adminRepo repository.AdminRepository, technicianRepo repository.TechnicianRepository, userRepo repository.UserRepository) {
	seeder.SeedRooms(roomRepo, 10)
	seeder.SeedAdmins(adminRepo, 5)
	seeder.SeedTechnicians(technicianRepo, adminRepo, 20)
	seeder.SeedUsers(userRepo, 40)
}
