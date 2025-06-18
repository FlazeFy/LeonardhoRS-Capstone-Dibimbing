package factory

import (
	"math/rand"
	"pelita/entity"

	"github.com/google/uuid"
)

var floors = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
var depts = []string{"IT", "Human Resource", "Finance & Risk Management", "Marketing", "Sales", "Planning & Transformation", "Network"}

func RandomRoomName() string {
	return "Room-" + uuid.New().String()[:8]
}

func RandomDept() string {
	return depts[rand.Intn(len(depts))]
}

func RandomFloor() string {
	return floors[rand.Intn(len(floors))]
}

func GenerateRoom() entity.Room {
	return entity.Room{
		Floor:    RandomFloor(),
		RoomName: RandomRoomName(),
		RoomDept: RandomDept(),
	}
}
