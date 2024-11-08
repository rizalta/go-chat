package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageRepo struct {
	db *redis.Client
}

func NewMessageRepo(db *redis.Client) *MessageRepo {
	return &MessageRepo{db}
}

func (r *MessageRepo) AddMessage(ctx context.Context, message domain.Message) error {
	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("message:%d", time.Now().UnixNano())
	_, err = r.db.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.Set(ctx, key, messageJson, 0).Err(); err != nil {
			return err
		}
		return p.ZAdd(ctx, "messages", redis.Z{
			Member: key,
			Score:  float64(message.TimeStamp.UnixNano()),
		}).Err()
	})
	return err
}

func (r *MessageRepo) GetAllMessages(ctx context.Context) ([]*domain.Message, error) {
	keys, err := r.db.ZRange(ctx, "messages", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*domain.Message
	for _, key := range keys {
		var message domain.Message
		messageJson, err := r.db.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		err = json.Unmarshal([]byte(messageJson), &message)
		if err != nil {
			continue
		}
		messages = append(messages, &message)
	}

	return messages, nil
}
