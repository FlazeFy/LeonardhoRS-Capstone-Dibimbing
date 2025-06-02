package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type RoomService interface {
	GetAllRoom() ([]entity.Room, error)
	GetRoomAssetByFloorAndRoomName(floor, roomName string) ([]entity.RoomAsset, error)
	GetRoomAssetShortByFloorAndRoomName(floor, roomName string) ([]entity.RoomAssetShort, error)
	Create(room *entity.Room) error
	UpdateById(room *entity.Room, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
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

func (s *roomService) GetRoomAssetByFloorAndRoomName(floor, roomName string) ([]entity.RoomAsset, error) {
	// Repo : Get Find Room Asset By Floor And Room Name
	roomAsset, err := s.roomRepo.FindRoomAssetByFloorAndRoomName(floor, roomName)
	if err != nil {
		return nil, err
	}
	if roomAsset == nil {
		return nil, errors.New("room not found")
	}

	return roomAsset, nil
}

func (s *roomService) GetRoomAssetShortByFloorAndRoomName(floor, roomName string) ([]entity.RoomAssetShort, error) {
	// Repo : Get Find Room Asset Short By Floor And Room Name
	roomAsset, err := s.roomRepo.FindRoomAssetShortByFloorAndRoomName(floor, roomName)
	if err != nil {
		return nil, err
	}
	if roomAsset == nil {
		return nil, errors.New("room not found")
	}

	return roomAsset, nil
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
	if err := s.roomRepo.Create(room); err != nil {
		return err
	}

	return nil
}

func (s *roomService) UpdateById(room *entity.Room, id uuid.UUID) error {
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
	is_exist, err := s.roomRepo.FindByRoomNameFloorAndId(room.RoomName, room.Floor, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("room already exist on the same floor")
	}

	// Repo : Update Room By Id
	if err := s.roomRepo.UpdateById(room, id); err != nil {
		return err
	}

	return nil
}

func (s *roomService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Room By Id
	err := s.roomRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}
