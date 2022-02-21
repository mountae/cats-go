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
	CreateUserServ(user models.User) (models.User, error)
	GenerateToken(username string, password string) (t string, rt string, err error)
	RefreshTokens(rt string) (nt, nrt string, err error)
}

func NewUserAuthService(r repository.Auth) *UserAuthService {
	return &UserAuthService{repository: r}
}

type JwtCustomClaims struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	jwt.StandardClaims
}

func (s *UserAuthService) CreateUserServ(user models.User) (models.User, error) {
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
		Name: user.Username,
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
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rfc)

	rt, err = refToken.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", "", errors.New("error during generate refresh token")
	}

	return t, rt, nil
}

func (s *UserAuthService) RefreshTokens(rt string) (nt string, nrt string, err error) {
	verifyResult, error := VerifyToken(rt)

	if verifyResult == nil {
		return "", "", error
	}
	ncl := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}

	ntoken := jwt.NewWithClaims(jwt.SigningMethodHS256, ncl)

	// Generate encoded token and send it as response.
	nt, err = ntoken.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", "", errors.New("error during generate new token")
	}

	nrfc := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}

	nrefToken := jwt.NewWithClaims(jwt.SigningMethodHS256, nrfc)

	nrt, err = nrefToken.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", "", errors.New("error during generate new refresh token")
	}
	return nt, nrt, nil
}

func VerifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")), nil
	})
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return token, nil
}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("SALT_FOR_GENERATE_PASSWORD"))))
}
