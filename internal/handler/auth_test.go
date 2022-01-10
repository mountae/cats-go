package handler

import (
	"CatsGo/internal/request"
	"CatsGo/internal/service"
	"CatsGo/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var userHandler *UserAuthHandler

func TestUserAuthHandler_SignUp(t *testing.T) {
	TestTable := []struct {
		name             string
		inputJson        string
		exceptBody       string
		exceptStatusCode int
	}{
		{
			name:             "OK",
			inputJson:        `{"name":"Steve Jobs","username":"Steve","password":"Stev13jb7"}`,
			exceptBody:       "1",
			exceptStatusCode: 200,
		},
		{
			name:             "name is missing",
			inputJson:        `{"username":"Steve","password":"Stev13_jb7"}`,
			exceptBody:       `{"id":0,"name":"","username":"","password":""}`,
			exceptStatusCode: http.StatusBadRequest,
		},
		{
			name:             "username is missing",
			inputJson:        `{"name":"Steve Jobs","password":"Stev13_jb7"}`,
			exceptBody:       `{"id":0,"name":"","username":"","password":""}`,
			exceptStatusCode: http.StatusBadRequest,
		},
		{
			name:             "password is missing",
			inputJson:        `{"name":"Steve Jobs","username":"Steve"}`,
			exceptBody:       `{"id":0,"name":"","username":"","password":""}`,
			exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(TestCase.inputJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			var srv service.Auth
			srv = mock.NewMockCatServ()
			userHandler = NewUserAuthHandler(srv)

			if assert.NoError(t, userHandler.SignUp(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
				assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}

func TestUserAuthHandler_SignIn(t *testing.T) {
	TestTable := []struct {
		name             string
		inputJson        string
		exceptBody       string
		exceptStatusCode int
	}{
		{
			name:             "OK",
			inputJson:        `{"username":"Steve","password":"Stev13_jb7"}`,
			exceptStatusCode: http.StatusOK,
		},
		{
			name:             "Username isn't in",
			inputJson:        `{"password":"Stev13_jb7"}`,
			exceptStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Password isn't in",
			inputJson:        `{"username":"Steve"}`,
			exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(TestCase.inputJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			var srv service.Auth
			srv = mock.NewMockCatServ()
			userHandler = NewUserAuthHandler(srv)

			if assert.NoError(t, userHandler.SignIn(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}
