// Package models contains structs that affect database entities
package models

import (
	"github.com/google/uuid"
)

// Cats contains all related data to cats in database
type Cats struct {
	ID   uuid.UUID `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name"`
}

// User contains all related data to user in database
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}
