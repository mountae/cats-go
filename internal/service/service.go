package service

import (
	"CatsGo/internal/models"
	"CatsGo/internal/repository"

	"github.com/google/uuid"
)

type CatService struct {
	repository repository.Repository
}

type Service interface {
	GetAllCatsServ() ([]*models.Cats, error)
	CreateCatsServ(cats models.Cats) (*models.Cats, error)
	GetCatServ(id uuid.UUID) *models.Cats
	UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error)
	DeleteCatServ(id uuid.UUID)
}

func NewCatService(rps repository.Repository) *CatService {
	return &CatService{repository: rps}
}

func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

func (s *CatService) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	return s.repository.CreateCat(cats)
}

func (s *CatService) GetCatServ(id uuid.UUID) *models.Cats {
	return s.repository.GetCat(id)
}

func (s *CatService) UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	return s.repository.UpdateCat(id, cats)
}

func (s *CatService) DeleteCatServ(id uuid.UUID) {
	s.repository.DeleteCat(id)
}
