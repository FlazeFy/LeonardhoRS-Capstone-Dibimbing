package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedTechnicians(repo repository.TechnicianRepository, adminRepo repository.AdminRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		admin, _ := adminRepo.FindOneRandom()
		technician := factory.GenerateTechnician()
		err := repo.Create(&technician, admin.ID)
		if err != nil {
			fmt.Printf("failed to seed technician %d: %v\n", i, err)
		}
	}
}
