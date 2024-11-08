package database

import (
	"context"
	"fmt"
	"go-chat/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageRepo struct {
	db *redis.Client
}

func (r *MessageRepo) AddMessage(ctx context.Context, message domain.Message) {
	r.db.Set(ctx, fmt.Sprintf("message:%d", time.Now().UnixNano()), message.Content, 0)
}

func (r *MessageRepo) GetAllMessages(ctx context.Context) {
}
