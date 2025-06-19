package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedRooms(repo repository.RoomRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		room := factory.GenerateRoom()
		err := repo.Create(&room)
		if err != nil {
			fmt.Printf("failed to seed room %d: %v\n", i, err)
		}
	}
}
