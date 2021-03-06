// Package mock services
package mock

import (
	"CatsGo/internal/models"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// CatServ is an autogenerated mock type for the Cats type
type CatServ struct {
	mock.Mock
}

// NewCatServ create new MockCatServ
func NewCatServ() *CatServ {
	return &CatServ{}
}

// GetAllCatsServ provides request for all cats
func (m *CatServ) GetAllCatsServ() ([]*models.Cats, error) {
	cat := models.Cats{
		ID:   uuid.New(),
		Name: "",
	}
	allcats := []*models.Cats{&cat}
	return allcats, nil
}

// CreateCatServ provides request for creating new cat
func (m *CatServ) CreateCatServ(cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

// GetCatServ provides request to get cat by 'id'
func (m *CatServ) GetCatServ(id uuid.UUID) (*models.Cats, error) {
	cat := models.Cats{
		ID:   uuid.New(),
		Name: "Steve Jobs",
	}
	return &cat, nil
}

// UpdateCatServ provides request to update cat
func (m *CatServ) UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

// DeleteCatServ provides request to delete cat
func (m *CatServ) DeleteCatServ(id uuid.UUID) (*models.Cats, error) {
	cat := models.Cats{
		ID:   uuid.New(),
		Name: "Steve Jobs",
	}
	return &cat, nil
}

// CreateUserServ provides request to create user
func (m *CatServ) CreateUserServ(user models.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	return user.ID, nil
}

// GenerateToken provides request to get a pair of tokens
func (m *CatServ) GenerateToken(username, password string) (t, rt string, err error) {
	t = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywibmFtZSI6InRlc3QxbiIsImV4cCI6MTY0MjI2MzQwOX0.dOEFgYBqu9Wt-I-F4iV7_hb-0AoyH_jI2lFFVauSUS8"
	rt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywibmFtZSI6IiIsImV4cCI6MTY0MjI3MzMwOX0.oVCA8L8rUbZG2DoJF2dH6ykx2u_e0rh6hgrz6ip9cU8"
	return t, rt, nil
}
