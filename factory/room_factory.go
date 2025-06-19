package factory

import (
	"pelita/entity"
	"pelita/utils"
	"strings"

	"github.com/google/uuid"
)

var floors = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
var depts = []string{"IT", "Human Resource", "Finance & Risk Management", "Marketing", "Sales", "Planning & Transformation", "Network"}

func RandomRoomName() string {
	return "Room-" + strings.ToUpper(uuid.New().String()[:4])
}

func GenerateRoom() entity.Room {
	return entity.Room{
		Floor:    utils.RandomPicker(floors),
		RoomName: RandomRoomName(),
		RoomDept: utils.RandomPicker(depts),
	}
}
