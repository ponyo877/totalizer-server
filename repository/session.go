package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ponyo877/totalizer-server/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db  *gorm.DB
	kvs *redis.Client
	sub *redis.PubSub
}

type Question struct {
	ID       string    `gorm:"column:beast_id, primaryKey"`
	RoomID   string    `gorm:"column:room_id"`
	Content  string    `gorm:"column:room_id"`
	Vote     int       `gorm:"column:vote"`
	CreateAt time.Time `gorm:"column:create_at"`
}

func (*Question) TableName() string {
	return "totalizer.question"
}

func NewSessionRepository(db *gorm.DB, kvs *redis.Client) *sessionRepository {
	return &sessionRepository{db, kvs, nil}
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

func (r *sessionRepository) OpenRoom(string) error {
	return nil
}

func (r *sessionRepository) SubscribeRoom(roomID string) error {
	r.sub = r.kvs.Subscribe(context.Background(), roomID)
	return nil
}

func (r *sessionRepository) IncrimentEnterCount(roomID string) (int, error) {
	key := fmt.Sprintf("room:%s:enter", roomID)
	value, err := r.kvs.Incr(context.Background(), key).Result()
	return int(value), err
}

func (r *sessionRepository) CreateQuestion(question *domain.Question) error {
	q := &Question{
		ID:       question.ID(),
		RoomID:   question.RoomID(),
		Content:  question.Content(),
		Vote:     question.Vote(),
		CreateAt: question.CreateAt(),
	}
	return r.db.Create(q).Error
}

func (r *sessionRepository) PublishQuestion(question *domain.Question) error {
	return r.kvs.Publish(context.Background(), question.RoomID(), question.Content()).Err()
}

func (r *sessionRepository) GetVoteCount(question string) (int, error) {
	key := fmt.Sprintf("question:%s:vote", question)
	countStr, err := r.kvs.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(countStr)
}

func (r *sessionRepository) GetEnterCount(roomID string) (int, error) {
	key := fmt.Sprintf("room:%s:enter", roomID)
	countStr, err := r.kvs.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(countStr)
}

func (r *sessionRepository) IncrimentVoteCount(roomID string, answer string) (int, error) {
	key := fmt.Sprintf("room:%s:vote:%s", roomID, answer)
	value, err := r.kvs.Incr(context.Background(), key).Result()
	return int(value), err
}

func (r *sessionRepository) PublishReady(roomID string) error {
	return r.kvs.Publish(context.Background(), roomID, "READY").Err()
}

func (r *sessionRepository) PublishResult(roomID string, questionID string) error {
	count, err := r.GetVoteCount(questionID)
	if err != nil {
		return err
	}
	return r.kvs.Publish(context.Background(), roomID, count).Err()
}

func (r *sessionRepository) UpdateQuestionVote(questionID string) error {
	count, err := r.GetVoteCount(questionID)
	if err != nil {
		return err
	}
	return r.db.Model(&Question{}).Where("id = ?", questionID).Update("vote", count).Error
}
