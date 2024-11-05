package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var (
	address  = os.Getenv("DB_ADDRESS")
	port     = os.Getenv("DB_PORT")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
)

func New() *redis.Client {
	num, err := strconv.Atoi(database)
	if err != nil {
		log.Fatalf(fmt.Sprintf("database incorrect %v", err))
	}

	fullAddress := fmt.Sprintf("%s:%s", address, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       num,
	})

	return rdb
}
