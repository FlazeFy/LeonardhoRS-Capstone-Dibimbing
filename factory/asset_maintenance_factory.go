package factory

import (
	"pelita/config"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func GenerateAssetMaintenance(assetPlacementId, technicianId uuid.UUID) entity.AssetMaintenance {
	desc := gofakeit.LoremIpsumSentence(10)

	// Random Hour
	startHour := gofakeit.Number(0, 18)
	duration := gofakeit.Number(1, 5)
	endHour := startHour + duration
	hourStart := entity.Time{Time: time.Date(0, 1, 1, startHour, 0, 0, 0, time.UTC)}
	hourEnd := entity.Time{Time: time.Date(0, 1, 1, endHour, 0, 0, 0, time.UTC)}

	return entity.AssetMaintenance{
		MaintenanceDay:       utils.RandomPicker(config.Days),
		MaintenanceHourStart: hourStart,
		MaintenanceHourEnd:   hourEnd,
		MaintenanceNotes:     &desc,
		AssetPlacementId:     assetPlacementId,
		MaintenanceBy:        technicianId,
	}
}
