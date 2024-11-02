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

func (s *service) AddUser(ctx context.Context, arg AddUserParams) (domain.User, error) {
	exists, err := s.isEmailRegistered(ctx, arg.Email)
	if err != nil {
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, errors.New("user: email already registered")
	}
	id := uuid.NewString()
	user := domain.User{
		ID:       id,
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
	}
	_, err = s.db.TxPipelined(ctx, func(p redis.Pipeliner) error {
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
		return domain.User{}, err
	}
	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	key := fmt.Sprintf("user:%s", id)
	userJson, err := s.db.Get(ctx, key).Result()
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

func (s *service) isEmailRegistered(ctx context.Context, email string) (bool, error) {
	exits, err := s.db.SIsMember(ctx, "registered_emails", email).Result()
	if err != nil {
		return false, err
	}
	return exits, nil
}
