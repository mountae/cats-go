package mock

import (
	"CatsGo/internal/models"
)

type MockCatServ struct {
}

func NewMockCatServ() *MockCatServ {
	return &MockCatServ{}
}

func (m *MockCatServ) GetAllCatsServ() ([]*models.Cats, error) {
	cat := models.Cats{
		ID:   0,
		Name: "",
	}
	allcats := []*models.Cats{&cat}
	return allcats, nil
}

func (m *MockCatServ) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

func (m *MockCatServ) GetCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Steve Jobs",
	}
	return &cat, nil
}

func (m *MockCatServ) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

func (m *MockCatServ) DeleteCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Steve Jobs",
	}
	return &cat, nil
}

func (m *MockCatServ) CreateUserServ(user models.User) (int, error) {
	user.ID = 1
	return user.ID, nil
}

func (m *MockCatServ) GenerateToken(username string, password string) (t string, rt string, err error) {
	// work token
	t = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywibmFtZSI6InRlc3QxbiIsImV4cCI6MTY0MjI2MzQwOX0.dOEFgYBqu9Wt-I-F4iV7_hb-0AoyH_jI2lFFVauSUS8"
	rt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywibmFtZSI6IiIsImV4cCI6MTY0MjI3MzMwOX0.oVCA8L8rUbZG2DoJF2dH6ykx2u_e0rh6hgrz6ip9cU8"
	return t, rt, nil
}
