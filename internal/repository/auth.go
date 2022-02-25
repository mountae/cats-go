package repository

import (
	"CatsGo/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth interface init
type Auth interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(username, password string) (models.User, error)
}

// CreateUser creates new user in pgdb
func (c *PostgresRepository) CreateUser(user models.User) (models.User, error) {
	var userData models.User

	id := uuid.New()
	row := c.conn.QueryRow(context.Background(), "INSERT INTO users (ID, Name, Username, Password) "+
		"VALUES ($1, $2, $3, $4) RETURNING id, name, username",
		id, user.Name, user.Username, user.Password)
	err := row.Scan(&userData.ID, &userData.Name, &userData.Username)
	if err != nil {
		log.Error(err)
		return userData, errors.New("error when adding to the database")
	}

	return userData, nil
}

// GetUser get user from pgdb
func (c *PostgresRepository) GetUser(username, password string) (models.User, error) {
	var user models.User

	err := c.conn.QueryRow(context.Background(), "SELECT id, name, username, password "+
		"FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Username, &user.Password)

	if err != nil {
		log.Error(err)
		return models.User{}, errors.New("user not in database")
	}

	if user.Password != password {
		return models.User{}, errors.New("incorrect password")
	}

	return user, nil
}

// CreateUser creates new user in mongodb
func (c *MongoRepository) CreateUser(user models.User) (models.User, error) {
	collection := c.client.Database("users").Collection("users")

	docs := []interface{}{
		bson.D{primitive.E{Key: "name", Value: user.Name}, {Key: "username", Value: user.Username},
			{Key: "password", Value: user.Password}},
	}

	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	return models.User{}, nil
}

// GetUser get user from mongodb
func (c *MongoRepository) GetUser(username, password string) (models.User, error) {
	return models.User{}, nil
}
