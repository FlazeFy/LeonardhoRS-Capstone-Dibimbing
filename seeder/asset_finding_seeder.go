package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func SeedAssetFindings(repo repository.AssetFindingRepository, assetPlacementRepo repository.AssetPlacementRepository, technicianRepo repository.TechnicianRepository, userRepo repository.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		assetPlacement, _ := assetPlacementRepo.FindOneRandom()
		assetFinding := factory.GenerateAssetFinding(assetPlacement.ID)
		randNumber := gofakeit.Number(0, 1)

		if randNumber == 0 {
			technician, _ := technicianRepo.FindOneRandom()
			err := repo.Create(&assetFinding, technician.ID, uuid.Nil)
			if err != nil {
				fmt.Printf("failed to seed asset finding %d for technician: %v\n", i, err)
			}
		} else {
			user, _ := userRepo.FindOneRandom()
			err := repo.Create(&assetFinding, uuid.Nil, user.ID)
			if err != nil {
				fmt.Printf("failed to seed asset finding %d for user: %v\n", i, err)
			}
		}
	}
}
