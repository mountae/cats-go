package repository

import (
	"CatsGo/internal/models"
	"context"
	"errors"
	"strconv"

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
	CreateCats(cats models.Cats) (*models.Cats, error)
	GetCat(id string) (*models.Cats, error)
	UpdateCat(id string, cats models.Cats) (*models.Cats, error)
	DeleteCat(id string) (*models.Cats, error)
}

func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{client: client}
}

func (c *PostgresRepository) GetAllCats() ([]*models.Cats, error) {

	var allcats []*models.Cats

	rows, err := c.conn.Query(context.Background(), "SELECT ID, name FROM cats")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		cats := models.Cats{
			ID:   0,
			Name: "",
		}

		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		cats.ID = values[0].(int32)
		cats.Name = values[1].(string)
		allcats = append(allcats, &cats)
	}

	return allcats, nil
}

func (c *PostgresRepository) CreateCats(cats models.Cats) (*models.Cats, error) {

	// Add new cat to DB
	commandTag, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", cats.ID, cats.Name)
	if err != nil {
		return &cats, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cats, errors.New("failed to create a cat")
	}

	return &cats, nil
}

func (c *PostgresRepository) GetCat(id string) (*models.Cats, error) {

	var cat models.Cats

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}

	// Get 'name'
	var name string
	err = c.conn.QueryRow(context.Background(), "SELECT name FROM cats WHERE id=$1", id).Scan(&name)
	if err != nil {
		return &cat, err
	}

	// Set params for models.Cats
	cat.ID = int32(idInt)
	cat.Name = name

	return &cat, nil
}

func (c *PostgresRepository) UpdateCat(id string, cats models.Cats) (*models.Cats, error) {

	// Conv id -> int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cats, err
	}

	// Refresh models.Cat
	cats.ID = int32(idInt)

	// Update DB
	commandTag, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", cats.Name, cats.ID)
	if err != nil {
		return &cats, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cats, errors.New("row isn't updated")
	}

	return &cats, nil
}

func (c *PostgresRepository) DeleteCat(id string) (*models.Cats, error) {

	var cat models.Cats

	// Conv id -> int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Refresh models.Cats
	cat.ID = int32(idInt)
	// Get 'name'
	var name string
	err = c.conn.QueryRow(context.Background(), "SELECT name FROM cats WHERE id=$1", id).Scan(&name)
	if err != nil {
		return &cat, err
	}
	cat.Name = name

	// Delete from DB
	commandTag, err := c.conn.Exec(context.Background(), "DELETE FROM cats WHERE id=$1", id)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("no such a row to delete")
	}

	return &cat, nil
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

func (c *MongoRepository) UpdateCat(id string, cats models.Cats) (*models.Cats, error) {
	// Conv id -> int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cats, err
	}

	// Refresh models.Cat
	cats.ID = int32(idInt)

	// Changing DB values
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	filter := bson.D{primitive.E{Key: "id", Value: cats.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: cats.Name}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return &cats, nil
}

func (c *MongoRepository) DeleteCat(id string) (*models.Cats, error) {
	var cat models.Cats

	// Conv id -> int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Refresh models.Cats
	cat.ID = int32(idInt)

	// Delete from DB
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	_, err = collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: idInt}})
	if err != nil {
		log.Fatal(err)
	}

	return &cat, nil
}
