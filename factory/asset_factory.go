package factory

import (
	"fmt"
	"pelita/config"
	"pelita/entity"
	"pelita/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateAsset() entity.Asset {
	desc := gofakeit.ProductDescription()
	merk := gofakeit.Company()
	price := fmt.Sprintf("%d", gofakeit.Number(1, 1000)*10000)

	return entity.Asset{
		AssetName:     gofakeit.ProductName(),
		AssetDesc:     &desc,
		AssetMerk:     &merk,
		AssetCategory: gofakeit.ProductCategory(),
		AssetPrice:    &price,
		AssetStatus:   utils.RandomPicker(config.AssetStatus),
	}
}
