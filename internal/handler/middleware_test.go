package handler

// import (
// 	"CatsGo/internal/service"
// 	"CatsGo/internal/service/mock"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// func TestUserAuthHandler_Restricted(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	var srv service.Auth
// 	srv = mock.NewMockCatServ()
// 	userHandler = NewUserAuthHandler(srv)

// 	if assert.NoError(t, userHandler.Restricted(c)) {
// 		// assert.Equal(t, 200, rec.Code)
// 		// assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))

// 	}
// }
