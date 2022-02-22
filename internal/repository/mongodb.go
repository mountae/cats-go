package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RequestMongo provides connection to MongoDB database
func RequestMongo() (*mongo.Client, context.CancelFunc) {
	const cont = 10

	if err := initConfig(); err != nil {
		log.Fatal("mongodb config files error")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("mongodb loading env variables error")
	}

	url := fmt.Sprintf("mongodb://%s:%s/",
		viper.GetString("mongodb.host"),
		viper.GetString("mongodb.port"))

	fmt.Println(url)

	url = os.Getenv("MONGODB_CONNSTRING")

	// comment out url when building
	// url = "mongodb://root:example@localhost:27017/"

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cont*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, cancel
}
