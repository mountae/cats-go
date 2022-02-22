package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// RequestDB provides connection to Postgres database
func RequestDB() *pgxpool.Pool {
	if err := initConfig(); err != nil {
		log.Fatal("postgres error with config files")
	}

	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		viper.GetString("db.pos"),
		viper.GetString("db.username"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbase"))

	fmt.Println(url)
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
