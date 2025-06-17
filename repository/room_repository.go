package repository

import (
	"errors"
	"fmt"
	"pelita/entity"
	"pelita/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room Interface
type RoomRepository interface {
	FindAll(pagination utils.Pagination) ([]entity.Room, int64, error)
	Create(room *entity.Room) error
	DeleteById(id uuid.UUID) error
	UpdateById(room *entity.Room, id uuid.UUID) error
	FindByRoomNameAndFloor(roomName, floor string) (*entity.Room, error)
	FindByRoomNameFloorAndId(roomName, floor string, id uuid.UUID) (*entity.Room, error)
	FindRoomAssetByFloorAndRoomName(floor, roomName string) ([]entity.RoomAsset, error)
	FindRoomAssetShortByFloorAndRoomName(floor, roomName string) ([]entity.RoomAssetShort, error)
}

// Room Struct
type roomRepository struct {
	db *gorm.DB
}

// Room Constructor
func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) FindAll(pagination utils.Pagination) ([]entity.Room, int64, error) {
	var total int64

	// Models
	var room []entity.Room

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	r.db.Model(&entity.Room{}).Count(&total)

	// Query
	err := r.db.Order("floor ASC").
		Order("room_name ASC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&room).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return room, total, nil
}

func (r *roomRepository) FindRoomAssetByFloorAndRoomName(floor, roomName string) ([]entity.RoomAsset, error) {
	// Models
	var roomAsset []entity.RoomAsset
	roomName = strings.ToLower(roomName)

	// Query
	var roomNameSelect string
	if roomName == "all" {
		roomNameSelect = "GROUP_CONCAT(DISTINCT room_name SEPARATOR ', ') as room_name"
	} else {
		roomNameSelect = "room_name"
	}

	query := r.db.Table("rooms").
		Select(fmt.Sprintf(`floor, %s, room_dept, asset_name, assets.asset_desc, SUM(asset_qty) as total_asset, asset_merk, asset_category`, roomNameSelect)).
		Joins("JOIN asset_placements ON asset_placements.room_id = rooms.id").
		Joins("JOIN assets ON assets.id = asset_placements.asset_id").
		Where("floor = ?", floor)

	if roomName != "all" {
		query = query.Where("room_name = ?", roomName)
	}

	result := query.Group("assets.id").
		Order("assets.asset_name ASC").
		Find(&roomAsset)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(roomAsset) == 0 {
		return nil, errors.New("room asset not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return roomAsset, nil
}

func (r *roomRepository) FindRoomAssetShortByFloorAndRoomName(floor, roomName string) ([]entity.RoomAssetShort, error) {
	// Models
	var roomAsset []entity.RoomAssetShort
	roomName = strings.ToLower(roomName)

	// Query
	var roomNameSelect string
	if roomName == "all" {
		roomNameSelect = "GROUP_CONCAT(DISTINCT room_name SEPARATOR ', ') as room_name"
	} else {
		roomNameSelect = "room_name"
	}

	query := r.db.Table("rooms").
		Select(fmt.Sprintf(`floor, %s, room_dept, asset_name, asset_category`, roomNameSelect)).
		Joins("JOIN asset_placements ON asset_placements.room_id = rooms.id").
		Joins("JOIN assets ON assets.id = asset_placements.asset_id").
		Where("floor = ?", floor)

	if roomName != "all" {
		query = query.Where("room_name = ?", roomName)
	}

	result := query.Group("assets.id").
		Order("assets.asset_name ASC").
		Find(&roomAsset)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(roomAsset) == 0 {
		return nil, errors.New("room asset not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return roomAsset, nil
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
