package repository

import (
	"context"
	"time"

	"github.com/ponyo877/totalizer-server/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db  *gorm.DB
	kvs *redis.Client
}

type Question struct {
	ID       string
	RoomID   string
	Content  string
	Vote     int
	CreateAt time.Time
}

func (*Question) TableName() string {
	return "totalizer.question"
}

func NewSessionRepository(db *gorm.DB, kvs *redis.Client) *sessionRepository {
	return &sessionRepository{db, kvs}
}

func (r *sessionRepository) Incriment(key string) (int, error) {
	value, err := r.kvs.Incr(context.Background(), key).Result()
	return int(value), err
}

func (r *sessionRepository) ListQuestion() (*domain.Question, error) {
	var questions []Question
	if err := r.db.Find(&questions).Error; err != nil {
		return nil, err
	}
	question := questions[2]
	return domain.NewQuestion(question.ID, question.RoomID, question.Content, question.Vote, question.CreateAt), nil
}
