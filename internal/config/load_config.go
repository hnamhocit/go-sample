package config

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sample/internal/database"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadConfig() (*database.Queries, context.Context) {
	LoadEnv()

	mdb, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	dao := database.New(mdb)
	ctx := context.Background()

	return dao, ctx
}
