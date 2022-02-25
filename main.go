package main

import (
	"CatsGo/internal/configs"
	"CatsGo/internal/handler"
	repo "CatsGo/internal/repository"
	"CatsGo/internal/request"
	"CatsGo/internal/service"
	"context"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "CatsGo/docs"
)

const (
	flag     = "postgres" // postgres / mongodb
	portEcho = ":8000"
)

// @title Cats Go
// @version 1.0
// @description This is a simple CRUD app for Go.

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	e := echo.New()
	e.Validator = &request.CustomValidator{Validator: validator.New()}

	// Configuration
	cfg := configs.Config{}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is Cats Go app!")
	})

	var ctx = context.TODO()
	var rps repo.Repository
	var rpsAuth repo.Auth
	if flag == "postgres" {
		// postgres connect
		conn, err := NewPgxPool(ctx, cfg)
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		rps = repo.NewPostgresRepository(conn)
		rpsAuth = repo.NewPostgresRepository(conn)
	} else if flag == "mongodb" {
		// mongodb connect
		client, err := NewMongoClient(ctx, cfg)
		if err != nil {
			log.Panic(err)
		}

		rps = repo.NewMongoRepository(client, cfg)
		rpsAuth = repo.NewMongoRepository(client, cfg)
	}

	var srv service.Service = service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)

	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCat)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	var srvAuth service.Auth = service.NewUserAuthService(rpsAuth, cfg)
	hndlrAuth := handler.NewUserAuthHandler(srvAuth)
	e.POST("/register", hndlrAuth.SignUp)
	e.POST("/login", hndlrAuth.SignIn)
	e.POST("/token", hndlrAuth.UpdateTokens)

	r := e.Group("/restrict")
	{
		config := middleware.JWTConfig{
			Claims:     new(service.JwtCustomClaims),
			SigningKey: []byte(cfg.KeyForSignatureJwt),
		}
		r.Use(middleware.JWTWithConfig(config))
		r.GET("", hndlrAuth.Restricted)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(portEcho))
}

// NewPgxPool provides connection with postgres database
func NewPgxPool(ctx context.Context, cfg configs.Config) (*pgxpool.Pool, error) {
	cfgErr := env.Parse(&cfg)
	if cfgErr != nil {
		log.Errorf("Unable to parse config: %v\n", cfgErr)
		return nil, fmt.Errorf("we can't parse config")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDBName)

	conn, err := pgxpool.Connect(ctx, url)
	fmt.Println(url) //TODO: fix binding cfg env vars
	if err != nil {
		log.Errorf("Unable to connect to postgres database: %v\n", err)
		return nil, fmt.Errorf("we can't connect to database")
	}
	return conn, nil
}

// NewMongoClient provides connection to mongodb database
func NewMongoClient(ctx context.Context, cfg configs.Config) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("we can't connect to database")
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("Unable to connect to mongodb database: %v\n", err)
		return nil, fmt.Errorf("we can't connect to database")
	}
	return client, nil
}
