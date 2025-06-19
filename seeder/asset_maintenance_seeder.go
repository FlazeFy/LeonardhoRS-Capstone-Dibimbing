package seeder

import (
	"fmt"
	"pelita/factory"
	"pelita/repository"
)

func SeedAssetMaintenances(repo repository.AssetMaintenanceRepository, adminRepo repository.AdminRepository, technicianRepo repository.TechnicianRepository, assetPlacementRepo repository.AssetPlacementRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	for i := 0; i < count; i++ {
		admin, _ := adminRepo.FindOneRandom()
		technician, _ := technicianRepo.FindOneRandom()
		assetPlacement, _ := assetPlacementRepo.FindOneRandom()
		assetMaintenance := factory.GenerateAssetMaintenance(assetPlacement.ID, technician.ID)

		err := repo.Create(&assetMaintenance, admin.ID)
		if err != nil {
			fmt.Printf("failed to seed asset maintenance %d: %v\n", i, err)
		}
	}
}
