package repository

import (
	"CatsGo/internal/models"
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostgresRepository struct {
	conn *pgxpool.Pool
}

type MongoRepository struct {
	client *mongo.Client
}

type Repository interface {
	GetAllCats() ([]*models.Cats, error)
	CreateCat(cats models.Cats) (*models.Cats, error)
	GetCat(id uuid.UUID) *models.Cats
	UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error)
	DeleteCat(id uuid.UUID)
}

func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{client: client}
}

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

func (c *PostgresRepository) CreateCat(cat models.Cats) (*models.Cats, error) {

	cat.ID = uuid.New()
	// Add new cat to DB
	commandTag, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", cat.ID, cat.Name)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("failed to create a cat")
	}

	return &cat, nil
}

func (c *PostgresRepository) GetCat(id uuid.UUID) *models.Cats {

	var cat models.Cats

	result := c.conn.QueryRow(context.Background(), "SELECT * FROM cats WHERE id=$1", id)

	err := result.Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil
	}

	return &cat
}

func (c *PostgresRepository) UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error) {

	// Update DB
	result, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", cats.Name, id)
	if err != nil {
		return &cats, err
	}
	if result.RowsAffected() != 1 {
		return &cats, errors.New("row isn't updated")
	}

	return &cats, nil
}

func (c *PostgresRepository) DeleteCat(id uuid.UUID) {

	// Delete from DB
	c.conn.Exec(context.Background(), "DELETE FROM cats WHERE id=$1", id)

}

func (c *MongoRepository) GetAllCats() ([]*models.Cats, error) {

	var allcats = []*models.Cats{}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	cur, currErr := collection.Find(context.TODO(), bson.D{})
	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(context.TODO())

	if err := cur.All(context.TODO(), &allcats); err != nil {
		panic(err)
	}
	return allcats, nil
}

func (c *MongoRepository) CreateCats(cats models.Cats) (*models.Cats, error) {

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

	docs := []interface{}{
		bson.D{primitive.E{Key: "id", Value: cats.ID}, {Key: "name", Value: cats.Name}},
	}

	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	return &cats, nil
}

func (c *MongoRepository) GetCat(id string) (*models.Cats, error) {
	var cat models.Cats

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: idInt}}).Decode(&cat)
	if err != nil {
		return &cat, err
	}
	return &cat, nil
}

func (c *MongoRepository) UpdateCat(id uuid.UUID, cats models.Cats) (*models.Cats, error) {

	// Changing DB values
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	filter := bson.D{primitive.E{Key: "id", Value: cats.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: cats.Name}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return &cats, nil
}

func (c *MongoRepository) DeleteCat(id string) (*models.Cats, error) {
	var cat models.Cats

	// Delete from DB
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	_, err := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Fatal(err)
	}

	return &cat, nil
}
