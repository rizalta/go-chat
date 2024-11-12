package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

const numLoad = 20

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

func (r *MessageRepo) GetMessages(
	ctx context.Context,
	page int64,
) ([]*domain.Message, bool, error) {
	start := page * numLoad
	end := start + numLoad - 1
	total, err := r.db.ZCard(ctx, "messages").Result()
	if err != nil {
		return nil, false, err
	}

	keys, err := r.db.ZRevRange(ctx, "messages", start, end).Result()
	if err != nil {
		return nil, false, err
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

	return messages, end < total-1, nil
}
