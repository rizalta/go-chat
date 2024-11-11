package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	address  = os.Getenv("DB_ADDRESS")
	port     = os.Getenv("DB_PORT")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
)

func New(ctx context.Context) (*redis.Client, error) {
	num, err := strconv.Atoi(database)
	if err != nil {
		log.Printf(fmt.Sprintf("database incorrect %v", err))
		return nil, err
	}

	fullAddress := fmt.Sprintf("%s:%s", address, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       num,
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Printf("database not connected: %v", err)
		return nil, err
	}

	return rdb, nil
}
