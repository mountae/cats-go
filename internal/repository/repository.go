// Package repository encapsulate work with database
package repository

import (
	"CatsGo/internal/configs"
	"CatsGo/internal/models"

	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostgresRepository init pgsql
type PostgresRepository struct {
	conn *pgxpool.Pool
}

// MongoRepository init mongodb
type MongoRepository struct {
	client *mongo.Client
	cfg    *configs.Config
}

// Repository contains methods for work with cats collection
type Repository interface {
	GetAllCats() ([]*models.Cats, error)
	CreateCat(cats models.Cats) (*models.Cats, error)
	GetCat(id uuid.UUID) *models.Cats
	UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error)
	DeleteCat(id uuid.UUID) error
}

// NewPostgresRepository creates new cats repository
func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

// NewMongoRepository creates new cats repository
func NewMongoRepository(client *mongo.Client, cfg configs.Config) *MongoRepository {
	return &MongoRepository{client: client, cfg: &cfg}
}

// GetAllCats provides request to get all cats from pgdb
func (c *PostgresRepository) GetAllCats() ([]*models.Cats, error) {
	var allcats []*models.Cats

	rows, err := c.conn.Query(context.Background(), "SELECT id, name FROM cats")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var cat models.Cats

		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}

		allcats = append(allcats, &cat)
	}
	return allcats, nil
}

// CreateCat provides request to create new cat in pgdb
func (c *PostgresRepository) CreateCat(cat models.Cats) (*models.Cats, error) {
	cat.ID = uuid.New()
	result, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", cat.ID, cat.Name)
	if err != nil {
		return &cat, err
	}
	if result.RowsAffected() != 1 {
		return &cat, errors.New("failed to create a cat")
	}
	return &cat, nil
}

// GetCat provides request to get cat by 'id' from pgdb
func (c *PostgresRepository) GetCat(id uuid.UUID) *models.Cats {
	var cat models.Cats

	result := c.conn.QueryRow(context.Background(), "SELECT * FROM cats WHERE id=$1", id)
	err := result.Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil
	}
	return &cat
}

// UpdateCat provides request to update cat by 'id' in pgdb
func (c *PostgresRepository) UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	result, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", cats.Name, id)
	if err != nil {
		return &cats, err
	}
	if result.RowsAffected() != 1 {
		return &cats, errors.New("row isn't updated")
	}
	return &cats, nil
}

// DeleteCat provides request to delete cat by 'id' from pgdb
func (c *PostgresRepository) DeleteCat(id uuid.UUID) error {
	_, err := c.conn.Exec(context.Background(), "DELETE FROM cats WHERE id=$1", id)
	if err != nil {
		return errors.New("error while deleting cat")
	}
	return nil
}

// GetAllCats provides request to get all cats from mongodb
func (c *MongoRepository) GetAllCats() ([]*models.Cats, error) {
	var allcats = []*models.Cats{}

	// collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	collection := c.client.Database(c.cfg.MongoDBName).Collection(c.cfg.MongoCollection)
	cur, currErr := collection.Find(context.TODO(), bson.D{})
	if currErr != nil {
		panic(currErr)
	}
	if err := cur.Close(context.TODO()); err != nil {
		panic(err)
	}
	if err := cur.All(context.TODO(), &allcats); err != nil {
		panic(err)
	}
	return allcats, nil
}

// CreateCat provides request to create cat in mongodb
func (c *MongoRepository) CreateCat(cats models.Cats) (*models.Cats, error) {
	// collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	collection := c.client.Database(c.cfg.MongoDBName).Collection(c.cfg.MongoCollection)
	docs := []interface{}{
		bson.D{primitive.E{Key: "id", Value: cats.ID}, {Key: "name", Value: cats.Name}},
	}
	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	return &cats, nil
}

// GetCat provides request to get cat by 'id' from mongodb
func (c *MongoRepository) GetCat(id uuid.UUID) *models.Cats {
	var cat models.Cats

	// collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	collection := c.client.Database(c.cfg.MongoDBName).Collection(c.cfg.MongoCollection)
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: cat.ID}}).Decode(&cat)
	if err != nil {
		return nil
	}
	return &cat
}

// UpdateCat provides request to update cat by 'id' in mongodb
func (c *MongoRepository) UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error) {
	// collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	collection := c.client.Database(c.cfg.MongoDBName).Collection(c.cfg.MongoCollection)
	filter := bson.D{primitive.E{Key: "id", Value: cats.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: cats.Name}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return &cats, nil
}

// DeleteCat provides request to delete cat by 'id' from mongodb
func (c *MongoRepository) DeleteCat(id uuid.UUID) error {
	// collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	collection := c.client.Database(c.cfg.MongoDBName).Collection(c.cfg.MongoCollection)
	_, err := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
