// Package service provides business logic of app
package service

import (
	"CatsGo/internal/models"
	"CatsGo/internal/repository"

	"github.com/google/uuid"
)

// CatService struct init
type CatService struct {
	repository repository.Repository
}

// Service contains logic for cats
type Service interface {
	GetAllCatsServ() ([]*models.Cats, error)
	CreateCatServ(cats models.Cats) (*models.Cats, error)
	GetCatServ(id uuid.UUID) *models.Cats
	UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error)
	DeleteCatServ(id uuid.UUID) error
}

// NewCatService creates new cats service
func NewCatService(rps repository.Repository) *CatService {
	return &CatService{repository: rps}
}

// GetAllCatsServ provides service for GetAllCats method
func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

// CreateCatServ provides service for CreateCat method
func (s *CatService) CreateCatServ(cats models.Cats) (*models.Cats, error) {
	return s.repository.CreateCat(cats)
}

// GetCatServ provides service for GetCat method
func (s *CatService) GetCatServ(id uuid.UUID) *models.Cats {
	return s.repository.GetCat(id)
}

// UpdateCatServ provides service for UpdateCat method
func (s *CatService) UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	return s.repository.UpdateCat(id, cats)
}

// DeleteCatServ provides service for DeleteCat method
func (s *CatService) DeleteCatServ(id uuid.UUID) error {
	return s.repository.DeleteCat(id)
}
