package factory

import (
	"pelita/config"
	"pelita/entity"
	"pelita/utils"
	"strings"

	"github.com/google/uuid"
)

func RandomRoomName() string {
	return "Room-" + strings.ToUpper(uuid.New().String()[:4])
}

func GenerateRoom() entity.Room {
	return entity.Room{
		Floor:    utils.RandomPicker(config.Floors),
		RoomName: RandomRoomName(),
		RoomDept: utils.RandomPicker(config.Departments),
	}
}
