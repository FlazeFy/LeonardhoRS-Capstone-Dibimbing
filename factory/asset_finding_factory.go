package factory

import (
	"pelita/config"
	"pelita/entity"
	"pelita/utils"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func GenerateAssetFinding(assetPlacementId uuid.UUID) entity.AssetFinding {
	desc := gofakeit.LoremIpsumSentence(10)

	return entity.AssetFinding{
		FindingCategory:  utils.RandomPicker(config.FindingCategories),
		FindingNotes:     desc,
		AssetPlacementId: assetPlacementId,
	}
}
