package handler

import (
	"CatsGo/internal/service"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Restricted provides access to hidden page for authorized users
// @Summary Restricted
// @Security ApiKeyAuth
// @Description example closed page
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Router /restrict [get]
func (h *UserAuthHandler) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name)
}
