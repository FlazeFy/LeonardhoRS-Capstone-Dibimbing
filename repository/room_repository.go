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
	FindByRoomNameAndFloor(roomName, floor string) (*entity.Room, error)
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

func (r *roomRepository) Create(room *entity.Room) error {
	// Query
	room.ID = uuid.New()

	return r.db.Create(room).Error
}
