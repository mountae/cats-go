// Package handler encapsulates work with http
package handler

import (
	"CatsGo/internal/models"
	"CatsGo/internal/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// CatHandler init
type CatHandler struct {
	src service.Service
}

// NewCatHandler creation
func NewCatHandler(srv service.Service) *CatHandler {
	return &CatHandler{src: srv}
}

// GetAllCats fetches all entities from cats collection
// @Summary GetAllCats
// @Tags Cats
// @Description collect all cats in array
// @Produce json
// @Success 200 {array} models.Cats
// @Router /cats [get]
func (h *CatHandler) GetAllCats(c echo.Context) error {
	allcats, err := h.src.GetAllCatsServ()
	if err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, allcats)
}

// CreateCat creates a new entity in cats collection
// @Summary CreateCat
// @Tags Cats
// @Description create cat
// @Accept json
// @Produce json
// @Param cats body models.Cats true "cats"
// @Success 201 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Router /cats [post]
func (h *CatHandler) CreateCat(c echo.Context) error {
	cats := new(models.Cats)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	cat, err := h.src.CreateCatServ(*cats)
	if err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, cat)
}

// GetCat fetches a single cat from cat collection by 'id'
// @Summary GetCat
// @Tags Cats
// @Description get cat by id
// @Accept json
// @Produce json
// @Param id path uuid.UUID true "id"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [get]
func (h *CatHandler) GetCat(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	cat, err := h.src.GetCatServ(id)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, cat)
}

// UpdateCat updates a single cat in cats collection by 'id'
// @Summary UpdateCat
// @Tags Cats
// @Description update cat by id
// @Accept json
// @Produce json
// @Param id path uuid.UUID true "id"
// @Param cats body models.Cats true "cats"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [put]
func (h *CatHandler) UpdateCat(c echo.Context) error {
	cats := new(models.Cats)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	id, _ := uuid.Parse(c.Param("id"))
	cat, err := h.src.UpdateCatServ(id, *cats)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, cat)
}

// DeleteCat deletes a single cat from cats collection by 'id'
// @Summary DeleteCat
// @Tags Cats
// @Description delete cat by id
// @Accept json
// @Produce json
// @Param id path uuid.UUID true "id"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCat(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	err := h.src.DeleteCatServ(id)
	if err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

// RequestCatID struct init
type RequestCatID struct {
	ID   uuid.UUID `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name"`
}
