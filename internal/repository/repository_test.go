package repository

import (
	"CatsGo/internal/models"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var db *pgxpool.Pool
var hostAndPort string

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort = resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	err = resource.Expire(120)
	if err != nil {
		log.Error(err)
		return
	} // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 180 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.Connect(context.Background(), databaseUrl)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// Migrations
func MyMigrate() {
	cmd := exec.Command("flyway", "-url=jdbc:postgresql://"+hostAndPort+"/dbname",
		"-user=user_name", "-password=secret", "migrate")
	cmd.Dir = "C:/Program Files/flyway-8.0.3"
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Number of elements in the table
const countOfCats = 2

// Repository initialization
var rps Repository
var rpsAuth Auth

func NewPostgresRepositoryTest(db *pgxpool.Pool) {
	rps = NewPostgresRepository(db)
	rpsAuth = NewPostgresRepository(db)
}

func TestInit(t *testing.T) {
	// Make Migrations & init Repository
	MyMigrate()
	NewPostgresRepositoryTest(db)
}

func TestPostgresRepository_GetAllCats(t *testing.T) {
	// GET all Cats
	allcats, err := rps.GetAllCats()

	// Type Mapping check
	typeAllcats := fmt.Sprintf("%T", allcats)
	var tr []*models.Cats
	TypeTrue := fmt.Sprintf("%T", tr)
	assert.Equal(t, typeAllcats, TypeTrue)

	// Check elements count at DB
	count := len(allcats)
	assert.Equal(t, count, countOfCats)

	// Errors check
	assert.Nil(t, err)
}

func TestPostgresRepository_CreateCats(t *testing.T) {
	// Input values
	catsTrue := models.Cats{
		ID:   3,
		Name: "Murzik3",
	}

	// Create Cat
	cat, err := rps.CreateCats(catsTrue)

	// Return values check
	assert.Equal(t, cat, &catsTrue)
	assert.Nil(t, err)

	// Adding values to DB check
	allcats, err := rps.GetAllCats()
	if err != nil {
		log.Fatal(err)
	}
	count := len(allcats)
	assert.Equal(t, count, countOfCats+1)
}

func TestPostgresRepository_GetCat(t *testing.T) {
	// Input values
	catTrue := models.Cats{
		ID:   3,
		Name: "Murzik3",
	}
	id := "3"

	// Get cat id
	cat, err := rps.GetCat(id)

	// Return values check
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}

func TestPostgresRepository_UpdateCat(t *testing.T) {
	// Input values
	id := "1"
	catTrue := models.Cats{
		ID:   1,
		Name: "Pushok1",
	}

	// Update cat
	cat, err := rps.UpdateCat(id, catTrue)

	// Return values check
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)

	// Adding values to DB check
	cat, err = rps.GetCat(id)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, cat, &catTrue)
}

func TestPostgresRepository_DeleteCat(t *testing.T) {
	// Input values
	catTrue := models.Cats{
		ID:   3,
		Name: "Murzik3",
	}
	id := "3"

	// Delete cat by id
	cat, err := rps.DeleteCat(id)

	// Return values check
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}

func TestPostgresRepository_CreateUser(t *testing.T) {
	TestTable := []struct {
		name        string
		user        models.User
		exceptError error
	}{
		{
			name: "OK",
			user: models.User{
				ID:       1,
				Name:     "Steve Jobs",
				Username: "Steve",
				Password: "Stev13_jb7",
			},
			exceptError: nil,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			id, err := rpsAuth.CreateUser(TestCase.user)

			// Return values check
			assert.Equal(t, TestCase.user.ID, id)
			assert.Equal(t, TestCase.exceptError, err)

			// Adding values to DB check
			var newUser models.User
			err = db.QueryRow(context.Background(), "SELECT id, name, username, password "+
				"FROM users WHERE username = $1", TestCase.user.Username).Scan(
				&newUser.ID, &newUser.Name, &newUser.Username, &newUser.Password)
			if err != nil {
				log.Fatal(err)
			}

			assert.Equal(t, TestCase.user, newUser, "User incorrect for saving")
		})
	}
}

func TestMongoRepository_GetUser(t *testing.T) {
	TestTable := []struct {
		name          string
		inputUsername string
		inputPassword string
		expectUser    models.User
		exceptError   error
	}{
		{
			name:          "OK",
			inputUsername: "Steve",
			inputPassword: "Stev13_jb7",
			expectUser: models.User{
				ID:       1,
				Name:     "Steve Jobs",
				Username: "Steve",
				Password: "Stev13_jb7",
			},
			exceptError: nil,
		},
		{
			name:          "user not in database",
			inputUsername: "Carl",
			inputPassword: "random",
			expectUser:    *new(models.User),
			exceptError:   errors.New("user not in database"),
		},
		{
			name:          "incorrect password",
			inputUsername: "Steve",
			inputPassword: "random",
			expectUser:    *new(models.User),
			exceptError:   errors.New("incorrect password"),
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			user, err := rpsAuth.GetUser(TestCase.inputUsername, TestCase.inputPassword)

			assert.Equal(t, TestCase.expectUser, user)
			assert.Equal(t, TestCase.exceptError, err)
		})
	}
}
