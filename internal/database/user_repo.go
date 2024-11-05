package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AddUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepo struct {
	db *redis.Client
}

func NewUserRepo(db *redis.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) AddUser(ctx context.Context, arg AddUserParams) (domain.User, error) {
	_, err := r.GetUserByEmail(ctx, arg.Email)
	if err == nil {
		return domain.User{}, errors.New("Email already exists")
	}
	id := uuid.NewString()
	user := domain.User{
		ID:       id,
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
	}
	_, err = r.db.TxPipelined(ctx, func(p redis.Pipeliner) error {
		userJson, err := json.Marshal(user)
		if err != nil {
			return nil
		}
		err = p.Set(ctx, fmt.Sprintf("user:%s", user.ID), userJson, 0).Err()
		if err != nil {
			return err
		}
		return p.Set(ctx, fmt.Sprintf("user_email:%s", user.Email), user.ID, 0).Err()
	})
	if err != nil {
		return domain.User{}, errors.New("Something went wrong")
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	key := fmt.Sprintf("user:%s", id)
	userJson, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	key := fmt.Sprintf("user_email:%s", email)
	id, err := r.db.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.User{}, errors.New("Email not registered")
		}
		return domain.User{}, errors.New("Something went wrong")
	}
	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, errors.New("Something went wrong")
	}
	return user, nil
}
