package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedUsers(repo repository.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		user := factory.GenerateUser()
		err := repo.Create(&user)
		if err != nil {
			fmt.Printf("failed to seed user %d: %v\n", i, err)
		}
	}
}
