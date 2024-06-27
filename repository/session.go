package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	kvs *redis.Client
}

func NewSessionRepository(kvs *redis.Client) *sessionRepository {
	return &sessionRepository{
		kvs,
	}
}

func (r *sessionRepository) Incriment(key string) (int, error) {
	value, err := r.kvs.Incr(context.Background(), key).Result()
	return int(value), err
}
