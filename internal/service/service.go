// Package service provides business logic of app
package service

import (
	"CatsGo/internal/models"
	"CatsGo/internal/repository"
	"github.com/labstack/gommon/log"

	"github.com/google/uuid"
)

// CatService interface of repository
type CatService struct {
	repository repository.Repository
	hash       repository.RedisRepository
}

// Service contains methods which get params from handler and sent them to repository
type Service interface {
	GetAllCatsServ() ([]*models.Cats, error)
	CreateCatServ(cats models.Cats) (*models.Cats, error)
	GetCatServ(id uuid.UUID) (*models.Cats, error)
	UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error)
	DeleteCatServ(id uuid.UUID) error
}

// NewCatService constructor
//func NewCatService(rps repository.Repository) *CatService {
//	return &CatService{repository: rps}
//}
func NewCatService(rps repository.Repository, hash repository.RedisRepository) *CatService {
	return &CatService{repository: rps, hash: hash}
}

// GetAllCatsServ called by handler and calls func in repository
func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

// CreateCatServ called by handler and calls func in repository
//func (s *CatService) CreateCatServ(cats models.Cats) (*models.Cats, error) {
//	return s.repository.CreateCat(cats)
//}
func (s *CatService) CreateCatServ(cats models.Cats) (*models.Cats, error) {
	err := s.hash.CreateCat(cats)
	if err != nil {
		log.Error(err)
	}
	return s.repository.CreateCat(cats)
}

// GetCatServ called by handler and calls func in repository
//func (s *CatService) GetCatServ(id uuid.UUID) (*models.Cats, error) {
//	return s.repository.GetCat(id)
//}
func (s *CatService) GetCatServ(id uuid.UUID) (*models.Cats, error) {
	cat, err := s.hash.GetCat(id)
	if err != nil {
		cat, err = s.repository.GetCat(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return cat, nil
}

// UpdateCatServ called by handler and calls func in repository
func (s *CatService) UpdateCatServ(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	return s.repository.UpdateCat(id, cats)
}

// DeleteCatServ called by handler and calls func in repository
//func (s *CatService) DeleteCatServ(id uuid.UUID) error {
//	return s.repository.DeleteCat(id)
//}
func (s *CatService) DeleteCatServ(id uuid.UUID) error {
	err := s.hash.DeleteCat(id)
	if err != nil {
		log.Error(err)
		return err
	}
	return s.repository.DeleteCat(id)
}
