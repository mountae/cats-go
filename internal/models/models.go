// Package models contains structs that affect database entities
package models

import (
	"github.com/google/uuid"
)

// Cats contains all related data to cats in database
type Cats struct {
	ID   uuid.UUID `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name" validate:"required,min=3"`
}

// User contains all related data to user in database
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required,min=3"`
	Username string    `json:"username" validate:"required,lowercase,min=4"`
	Password string    `json:"password" validate:"required,max=20,min=6"`
}
