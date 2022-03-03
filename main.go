package main

import (
	"CatsGo/internal/configs"
	"CatsGo/internal/handler"
	repo "CatsGo/internal/repository"
	"CatsGo/internal/request"
	"CatsGo/internal/service"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/go-redis/redis/v8"

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
	dir      = "files/media/"
)

// NewPgxPool provides connection with postgres database
func NewPgxPool(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDBName)
	conn, cfgErr := pgxpool.Connect(ctx, url)
	if cfgErr != nil {
		log.Errorf("unable to connect to postgres database: %v\n", cfgErr)
		return nil, fmt.Errorf("we can't connect to postgres database")
	}
	return conn, nil
}

// NewMongoClient provides connection with mongo database
func NewMongoClient(ctx context.Context, cfg *configs.Config) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("we can't setup connection with mongo database")
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("unable to connect to mongo database: %v\n", err)
		return nil, fmt.Errorf("we can't connect to mongo database")
	}
	return client, nil
}

// NewRedisClient provides connection with redis
func NewRedisClient(cfg *configs.Config) (*redis.Client, error) {
	rHostPort := cfg.RedisHost + ":" + cfg.RedisPort
	rdb := redis.NewClient(&redis.Options{
		Addr:     rHostPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb, nil
}

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
	cfg := &configs.Config{}
	opts := &env.Options{}
	if err := env.Parse(cfg, *opts); err != nil {
		log.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is Cats Go app!")
	})

	var (
		rps     repo.Repository
		rpsAuth repo.Auth
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

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

	// redis connect
	rdb, err := NewRedisClient(cfg)
	if err != nil {
		log.Panic(err)
	}
	rds := repo.NewRedisRepository(rdb)

	var srv service.Service = service.NewCatService(rps, *rds)
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

	// Download file
	e.GET("/download", func(c echo.Context) error {
		return c.File("files/template/download.html")
	})
	e.GET("/download/file", func(c echo.Context) error {
		return c.Attachment("files/media/gopl.jpg", "new-gopl.jpg")
	})

	// Upload file
	e.GET("/upload", func(c echo.Context) error {
		return c.File("files/template/upload.html")
	})
	e.POST("/upload", func(c echo.Context) error {
		name := c.FormValue("name")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			return err
		}
		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return err
		}
		defer func(src multipart.File) {
			err = src.Close()
		}(src)

		// Destination
		dst, err := os.Create(dir + file.Filename)
		if err != nil {
			log.Error("error while creating file")
			return err
		}
		defer func(dst *os.File) {
			err = dst.Close()
			if err != nil {
				log.Error("error with destination close file")
			}
		}(dst)

		// Copy the uploaded file to the created file on the filesystem
		if _, err = io.Copy(dst, src); err != nil {
			log.Error("error while copy file")
			return err
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with field name='%s' and size='%d' bytes.</p>", file.Filename, name, file.Size))
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(portEcho))
}
