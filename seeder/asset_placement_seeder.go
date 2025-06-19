package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedAssetPlacements(repo repository.AssetPlacementRepository, adminRepo repository.AdminRepository, roomRepo repository.RoomRepository, assetRepo repository.AssetRepository, technicianRepo repository.TechnicianRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		admin, _ := adminRepo.FindOneRandom()
		technician, _ := technicianRepo.FindOneRandom()
		room, _ := roomRepo.FindOneRandom()
		asset, _ := assetRepo.FindOneRandom()
		assetPlacement := factory.GenerateAssetPlacement(asset.ID, room.ID, technician.ID)

		err := repo.Create(&assetPlacement, admin.ID)
		if err != nil {
			fmt.Printf("failed to seed asset placement %d: %v\n", i, err)
		}
	}
}
