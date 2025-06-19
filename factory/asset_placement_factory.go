package factory

import (
	"pelita/entity"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func GenerateAssetPlacement(assetId, roomId, technicianId uuid.UUID) entity.AssetPlacement {
	desc := gofakeit.LoremIpsumSentence(10)

	return entity.AssetPlacement{
		AssetQty:   gofakeit.Number(1, 30),
		AssetDesc:  &desc,
		AssetId:    assetId,
		RoomId:     roomId,
		AssetOwner: technicianId,
	}
}
