package repository

import (
	"errors"
	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	FindAll() ([]entity.Room, error)
	Create(room *entity.Room) error
	DeleteById(id uuid.UUID) error
	UpdateById(room *entity.Room, id uuid.UUID) error
	FindByRoomNameAndFloor(roomName, floor string) (*entity.Room, error)
	FindByRoomNameFloorAndId(roomName, floor string, id uuid.UUID) (*entity.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) FindAll() ([]entity.Room, error) {
	// Models
	var room []entity.Room

	// Query
	err := r.db.Find(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return room, err
}

func (r *roomRepository) FindByRoomNameAndFloor(roomName, floor string) (*entity.Room, error) {
	// Models
	var room entity.Room

	// Query
	err := r.db.Where("room_name = ? AND floor = ?", roomName, floor).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &room, err
}

func (r *roomRepository) FindByRoomNameFloorAndId(roomName, floor string, id uuid.UUID) (*entity.Room, error) {
	// Models
	var room entity.Room

	// Query
	err := r.db.Where("room_name = ? AND floor = ? AND id != ?", roomName, floor, id).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &room, err
}

func (r *roomRepository) Create(room *entity.Room) error {
	room.ID = uuid.New()

	// Query
	return r.db.Create(room).Error
}

func (r *roomRepository) UpdateById(room *entity.Room, id uuid.UUID) error {
	// Query : Check Old Room
	var existingRoom entity.Room
	if err := r.db.First(&existingRoom, "id = ?", id).Error; err != nil {
		return err
	}

	// Query : Update
	room.ID = id
	room.CreatedAt = existingRoom.CreatedAt

	if err := r.db.Save(&room).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) DeleteById(id uuid.UUID) error {
	// Models
	var room entity.Room

	// Query
	err := r.db.Unscoped().Where("id = ?", id).Delete(&room).Error
	if err != nil {
		return err
	}

	return nil
}
