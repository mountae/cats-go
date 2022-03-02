package repository

import (
	"CatsGo/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RedisRepository provides a connection with redis
type RedisRepository struct {
	rdb *redis.Client
}

// NewRedisRepository is constructor
func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func (c *RedisRepository) CreateCat(cat models.Cats) error {
	ctx := context.TODO()
	cat.ID = uuid.New()
	catID := cat.ID.String()

	err := c.rdb.Set(ctx, catID, cat.Name, 0).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *RedisRepository) GetCat(id uuid.UUID) (*models.Cats, error) {
	ctx := context.TODO()
	//id, err := uuid.Parse(echo.Context.Param(ctx, "id"))
	cat, err := c.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fmt.Println(cat)
	return &models.Cats{}, nil
}

func (c *RedisRepository) DeleteCat(id uuid.UUID) error {
	ctx := context.TODO()

	_, err := c.rdb.Del(ctx, id.String()).Result()
	if err != nil {
		log.Error("redis error while deleting a cat")
		return err
	}
	return nil
}
