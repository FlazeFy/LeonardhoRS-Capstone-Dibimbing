package routes

import (
	"pelita/repository"
	"pelita/seeder"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, roomRepo repository.RoomRepository) {
	seeder.SeedRooms(roomRepo, 10)
}
