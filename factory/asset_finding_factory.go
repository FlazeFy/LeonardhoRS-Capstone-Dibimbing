package factory

import (
	"pelita/entity"
	"pelita/utils"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

var findingCategories = []string{"broken", "missing", "upgrade", "feedback"}

func GenerateAssetFinding(assetPlacementId uuid.UUID) entity.AssetFinding {
	desc := gofakeit.LoremIpsumSentence(10)

	return entity.AssetFinding{
		FindingCategory:  utils.RandomPicker(findingCategories),
		FindingNotes:     desc,
		AssetPlacementId: assetPlacementId,
	}
}
