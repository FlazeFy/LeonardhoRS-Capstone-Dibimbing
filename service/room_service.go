package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
)

type RoomService interface {
	GetAllRoom() ([]entity.Room, error)
	Create(room *entity.Room) error
}

type roomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (s *roomService) GetAllRoom() ([]entity.Room, error) {
	// Repo : Get All Room
	room, err := s.roomRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (s *roomService) Create(room *entity.Room) error {
	// Validator
	if room.RoomName == "" {
		return errors.New("room name is required")
	}
	if room.RoomDept == "" {
		return errors.New("room dept is required")
	}
	if room.Floor == "" {
		return errors.New("floor is required")
	}

	// Repo : Get Room by Room Name & Floor
	is_exist, err := s.roomRepo.FindByRoomNameAndFloor(room.RoomName, room.Floor)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("room already exist on the same floor")
	}

	// Repo : Create Room
	s.roomRepo.Create(room)

	return nil
}
