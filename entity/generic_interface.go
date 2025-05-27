package entity

import "github.com/google/uuid"

type Account interface {
	GetID() uuid.UUID
	GetPassword() string
}
