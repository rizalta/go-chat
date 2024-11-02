package database

import (
	"context"
	"fmt"
	"go-chat/internal/domain"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	AddUser(ctx context.Context, arg AddUserParams) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	Close()
}

type service struct {
	db *redis.Client
}

var (
	address  = os.Getenv("DB_ADDRESS")
	port     = os.Getenv("DB_PORT")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
)

func New() Service {
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
	s := &service{db: rdb}

	return s
}

func (s *service) Close() {
	s.db.Close()
}
