// Package repository encapsulate work with redis database
package repository

import (
	"CatsGo/internal/models"
	"context"

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

// CreateCat provides request to create new cat in redis database
func (c *RedisRepository) CreateCat(cat models.Cats) error {
	ctx := context.TODO()
	cat.ID = uuid.New()
	catID := cat.ID.String()
	args := catID + ":" + cat.Name

	err := c.rdb.Set(ctx, catID, args, 0).Err()
	if err != nil {
		log.Error("redis error while creating a cat")
		return err
	}
	return nil
}

// GetCat provides request to get cat by 'id' from redis database
func (c *RedisRepository) GetCat(id uuid.UUID) (*models.Cats, error) {
	ctx := context.TODO()
	catID := id.String()
	val, err := c.rdb.Get(ctx, catID).Result()
	if err != nil {
		log.Error("redis error no such a cat in database")
		return nil, err
	}
	return &models.Cats{ID: id, Name: val}, nil
}

// DeleteCat provides request to delete cat by 'id' from redis database
func (c *RedisRepository) DeleteCat(id uuid.UUID) error {
	ctx := context.TODO()

	_, err := c.rdb.Del(ctx, id.String()).Result()
	if err != nil {
		log.Error("redis error while deleting a cat")
		return err
	}
	return nil
}
