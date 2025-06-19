package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedAssets(repo repository.AssetRepository, adminRepo repository.AdminRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		admin, _ := adminRepo.FindOneRandom()
		asset := factory.GenerateAsset()
		err := repo.Create(&asset, admin.ID)
		if err != nil {
			fmt.Printf("failed to seed asset %d: %v\n", i, err)
		}
	}
}
