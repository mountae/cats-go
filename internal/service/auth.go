package service

import (
	"CatsGo/internal/models"
	"CatsGo/internal/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UserAuthService struct {
	repository repository.Auth
}

type Auth interface {
	CreateUserServ(user models.User) (int, error)
	GenerateToken(username string, password string) (t string, rt string, err error)
}

func NewUserAuthService(r repository.Auth) *UserAuthService {
	return &UserAuthService{repository: r}
}

type JwtCustomClaims struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	jwt.StandardClaims
}

func (s *UserAuthService) CreateUserServ(user models.User) (int, error) {
	user.Password = generatePassword(user.Password)
	return s.repository.CreateUser(user)
}

func (s *UserAuthService) GenerateToken(username string, password string) (t string, rt string, err error) {
	user, err := s.repository.GetUser(username, generatePassword(password))
	if err != nil {
		return "", "", errors.New("error with generate token in repository")
	}

	ac := &JwtCustomClaims{
		ID:   user.ID,
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)

	// Generate encoded token and send it as response.
	t, err = token.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", "", errors.New("error during generate token")
	}

	rfc := &JwtCustomClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rfc)

	rt, err = refToken.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", "", errors.New("error during generate refresh token")
	}

	return t, rt, nil
}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("SALT_FOR_GENERATE_PASSWORD"))))
}
