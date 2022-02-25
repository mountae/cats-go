// Package handler provides work with requests
package handler

import (
	"CatsGo/internal/models"
	"CatsGo/internal/service"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// UserAuthHandler init
type UserAuthHandler struct {
	src service.Auth
}

// NewUserAuthHandler creation
func NewUserAuthHandler(src service.Auth) *UserAuthHandler {
	return &UserAuthHandler{src: src}
}

// SignUp provides logic for register user
// @Summary SignUp
// @Tags auth
// @Description decode params and send it in service for create account
// @Accept json
// @Produce json
// @Param user body models.User true "user"
// @Success 200 {integer} integer
// @Failure 400 {object} models.User
// @Failure 500 {object} models.User
// @Router /register [post]
func (h *UserAuthHandler) SignUp(c echo.Context) error {
	var user models.User

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	if err = c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	id, err := h.src.CreateUserServ(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, id)
}

// SignInInput struct init
type SignInInput struct {
	Username string `json:"username" validate:"required,lowercase,min=4"`
	Password string `json:"password" validate:"required,max=20,min=6"`
}

// TokenResponse struct init
type TokenResponse struct {
	AccessToken  string `json:"accessToken" validate:"required,jwt"`
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}

// RefreshTokenRequest struct init
type RefreshTokenRequest struct {
	Token string `json:"Token" validate:"required,jwt"`
}

// SignIn provides logic for login user
// @Summary SignIn
// @Tags auth
// @Description decode params and send them in service for generate token
// @Accept json
// @Produce json
// @Param input body SignInInput true "input"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} models.User
// @Failure 500 {object} models.User
// @Router /login [post]
func (h *UserAuthHandler) SignIn(c echo.Context) error {
	var input SignInInput

	err := json.NewDecoder(c.Request().Body).Decode(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	if err = c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	token, refToken, err := h.src.GenerateToken(input.Username, input.Password)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	a := TokenResponse{AccessToken: token, RefreshToken: refToken}
	return c.JSON(http.StatusOK, a)
}

// UpdateTokens provides logic for update users tokens
// @Summary UpdateTokens
// @Tags auth
// @Description update access and refresh token pair
// @Accept json
// @Produce json
// @Param t_input body TokenRespone true "t_input"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} models.User
// @Failure 500 {object} models.User
// @Router /token [post]
func (h *UserAuthHandler) UpdateTokens(c echo.Context) error {
	var tInput RefreshTokenRequest

	err := json.NewDecoder(c.Request().Body).Decode(&tInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}
	if err = c.Validate(tInput); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}
	ntoken, nrefToken, err := h.src.RefreshTokens(tInput.Token)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	b := TokenResponse{AccessToken: ntoken, RefreshToken: nrefToken}
	return c.JSON(http.StatusOK, b)
}
