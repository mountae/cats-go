package repository

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RequestMongo() (*mongo.Client, context.CancelFunc) {
	if err := initConfig(); err != nil {
		log.Fatal("mongodb config files error")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("mongodb loading env variables error")
	}

	//url := fmt.Sprintf("mongodb://%s:%s/",
	//	viper.GetString("mongodb.host"),
	//	viper.GetString("mongodb.port"))
	//
	//fmt.Println(url)

	url := os.Getenv("MONGODB_CONNSTRING")

	// comment out url when building
	//url = "mongodb://root:example@localhost:27017/"

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, cancel
}
