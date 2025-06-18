package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedAdmins(repo repository.AdminRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		admin := factory.GenerateAdmin()
		err := repo.Create(&admin)
		if err != nil {
			fmt.Printf("failed to seed admin %d: %v\n", i, err)
		}
	}
}
