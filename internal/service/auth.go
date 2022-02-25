// Package service provides logic for authorization
package service

import (
	"CatsGo/internal/configs"
	"CatsGo/internal/models"
	"CatsGo/internal/repository"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

// UserAuthService implements an interface of Auth from repository
type UserAuthService struct {
	repository repository.Auth
	cfg        *configs.Config
}

// Auth contains methods for auth cases
type Auth interface {
	CreateUserServ(user models.User) (models.User, error)
	GenerateToken(username string, password string) (t string, rt string, err error)
	RefreshTokens(rt string) (nt, nrt string, err error)
}

// NewUserAuthService is a constructor
func NewUserAuthService(r repository.Auth, cfg configs.Config) *UserAuthService {
	return &UserAuthService{repository: r, cfg: &cfg}
}

// JwtCustomClaims expands the jwt.StandardClaims
type JwtCustomClaims struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	jwt.StandardClaims
}

// CreateUserServ provides new service for user
func (s *UserAuthService) CreateUserServ(user models.User) (models.User, error) {
	user.Password = generatePassword(user.Password, s.cfg)
	return s.repository.CreateUser(user)
}

// GenerateToken func creates a pair of jwt tokens
func (s *UserAuthService) GenerateToken(username, password string) (t, rt string, err error) {
	const (
		att = 15
		rtt = 1
	)

	user, err := s.repository.GetUser(username, generatePassword(password, s.cfg))
	if err != nil {
		log.Error("error with generate token in repository")
		return "", "", err
	}

	ac := &JwtCustomClaims{
		ID:   user.ID,
		Name: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * att).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)

	// Generate encoded token and send it as response.
	t, err = token.SignedString([]byte(s.cfg.KeyForSignatureJwt))
	if err != nil {
		log.Error("error during generate token")
		return "", "", err
	}

	rfc := &JwtCustomClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * rtt).Unix(),
		},
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rfc)

	rt, err = refToken.SignedString([]byte(s.cfg.KeyForSignatureJwt))
	if err != nil {
		log.Error("error during generate refresh token")
		return "", "", err
	}

	return t, rt, nil
}

// RefreshTokens func provides force update a pair of tokens
func (s *UserAuthService) RefreshTokens(rt string) (nt, nrt string, err error) {
	const (
		natt = 30
		nrtt = 3
	)
	verifyResult, err := VerifyToken(rt, s.cfg) // s.cfg hz

	if verifyResult == nil {
		log.Error("Token not verified")
		return "", "", err
	}
	ncl := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * natt).Unix(),
		},
	}

	ntoken := jwt.NewWithClaims(jwt.SigningMethodHS256, ncl)

	nt, err = ntoken.SignedString([]byte(s.cfg.KeyForSignatureJwt))
	if err != nil {
		log.Error("error during generate new token")
		return "", "", err
	}

	nrfc := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * nrtt).Unix(),
		},
	}

	nrefToken := jwt.NewWithClaims(jwt.SigningMethodHS256, nrfc)

	nrt, err = nrefToken.SignedString([]byte(s.cfg.KeyForSignatureJwt))
	if err != nil {
		log.Error("error during generate new refresh token")
		return "", "", err
	}
	return nt, nrt, nil
}

// VerifyToken func does validation for entered tokens
func VerifyToken(t string, cfg *configs.Config) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.KeyForSignatureJwt), nil
	})
	if _, ok := token.Claims.(jwt.StandardClaims); !ok && !token.Valid {
		return nil, err
	}

	if err != nil {
		log.Error("error while verify token")
		return nil, err
	}

	return token, nil
}

func generatePassword(password string, cfg *configs.Config) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	// return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("SALT_FOR_GENERATE_PASSWORD"))))
	return fmt.Sprintf("%x", hash.Sum([]byte(cfg.Salt)))
}
