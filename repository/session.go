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
	ID        string    `gorm:"column:id"`
	RoomID    string    `gorm:"column:room_id"`
	Content   string    `gorm:"column:content"`
	VoteCount int       `gorm:"column:vote_count"`
	YesCount  int       `gorm:"column:yes_count"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (*Question) TableName() string {
	return "totalizer.question"
}

func NewSessionRepository(db *gorm.DB, kvs *redis.Client) *sessionRepository {
	return &sessionRepository{db, kvs, nil}
}

func (r *sessionRepository) Incriment(key string) (int, error) {
	value, err := r.kvs.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return int(value), err
}

func (r *sessionRepository) ListQuestion() (*domain.Question, error) {
	var questions []Question
	if err := r.db.Find(&questions).Error; err != nil {
		return nil, err
	}
	question := questions[2]
	return domain.NewQuestion(question.ID, question.RoomID, question.Content, question.VoteCount, question.YesCount, question.CreatedAt), nil
}

func (r *sessionRepository) SubscribeRoom(roomID string) *chan string {
	r.sub = r.kvs.Subscribe(context.Background(), roomID)
	ch := make(chan string)
	go func() {
		subCh := r.sub.Channel()
		for msg := range subCh {
			ch <- msg.Payload
		}
	}()
	return &ch
}

func (r *sessionRepository) IncrimentEnterCount(roomID string) (int, error) {
	key := fmt.Sprintf("room:%s:enter", roomID)
	value, err := r.kvs.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return int(value), err
}

func (r *sessionRepository) CreateQuestion(question *domain.Question) error {
	q := &Question{
		ID:        question.ID(),
		RoomID:    question.RoomID(),
		Content:   question.Content(),
		VoteCount: question.VoteCount(),
		YesCount:  question.YesCount(),
		CreatedAt: question.CreatedAt(),
	}
	return r.db.Create(q).Error
}

func (r *sessionRepository) PublishQuestion(question *domain.Question) error {
	ans, err := domain.NewAskAnswer(question.ID(), question.Content())
	if err != nil {
		return err
	}
	answer, err := ans.String()
	if err != nil {
		return err
	}
	return r.kvs.Publish(context.Background(), question.RoomID(), answer).Err()
}

func (r *sessionRepository) GetVoteCount(questionID string) (int, error) {
	key := fmt.Sprintf("question:%s:vote", questionID)
	countStr, err := r.kvs.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(countStr)
}

func (r *sessionRepository) GetAnswerCount(questionID string, answer string) (int, error) {
	key := fmt.Sprintf("question:%s:vote:%s", questionID, answer)
	countStr, err := r.kvs.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	}
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

func (r *sessionRepository) IncrimentVoteCount(questionID string, answer string) (int, error) {
	key := fmt.Sprintf("question:%s:vote:%s", questionID, answer)
	value, err := r.kvs.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	key = fmt.Sprintf("question:%s:vote", questionID)
	if _, err := r.kvs.Incr(context.Background(), key).Result(); err != nil {
		return 0, err
	}
	return int(value), err
}

func (r *sessionRepository) PublishReady(roomID string) error {
	ans, err := domain.NewReadyAnswer()
	if err != nil {
		return err
	}
	answer, err := ans.String()
	if err != nil {
		return err
	}
	return r.kvs.Publish(context.Background(), roomID, answer).Err()
}

func (r *sessionRepository) PublishResult(roomID string, yesCount, enterCount int) error {
	ans, err := domain.NewResultAnswer(yesCount, enterCount)
	if err != nil {
		return err
	}
	answer, err := ans.String()
	if err != nil {
		return err
	}
	return r.kvs.Publish(context.Background(), roomID, answer).Err()
}

func (r *sessionRepository) UpdateQuestionVote(questionID string) error {
	vc, err := r.GetVoteCount(questionID)
	if err != nil {
		return err
	}
	yc, err := r.GetAnswerCount(questionID, "yes")
	if err != nil {
		return err
	}
	return r.db.Model(&Question{}).Updates(Question{ID: questionID, VoteCount: vc, YesCount: yc, UpdatedAt: time.Now()}).Error
}

func (r *sessionRepository) PublishEnter(roomID string, enterCount int) error {
	ans, err := domain.NewEnterAnswer(enterCount)
	if err != nil {
		return err
	}
	answer, err := ans.String()
	if err != nil {
		return err
	}
	return r.kvs.Publish(context.Background(), roomID, answer).Err()
}

func (r *sessionRepository) StoreRoomStatus(roomID string, status domain.RoomStatus) error {
	key := fmt.Sprintf("status:%s", roomID)
	return r.kvs.Set(context.Background(), key, int(status), 0).Err()
}

func (r *sessionRepository) GetRoomStatus(roomID string) (*domain.Status, error) {
	key := fmt.Sprintf("status:%s", roomID)
	statusStr, err := r.kvs.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, err
	}
	return domain.NewStatus(domain.RoomStatus(statusInt)), nil
}

func (r *sessionRepository) GetLatestQuestion(roomID string) (*domain.Question, error) {
	var q Question
	if err := r.db.Where("room_id = ?", roomID).Order("created_at desc").Limit(1).Find(&q).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewQuestion(q.ID, q.RoomID, q.Content, q.VoteCount, q.YesCount, q.CreatedAt), nil
}

func (r *sessionRepository) GetRoomIDByRoomNumber(roomNumber string) (string, bool, error) {
	key := fmt.Sprintf("room:%s", roomNumber)
	roomID, err := r.kvs.Get(context.Background(), key).Result()
	switch err {
	case nil:
		return roomID, false, nil
	case redis.Nil:
		return "", true, nil
	default:
		return "", false, err
	}
}

func (r *sessionRepository) SetRoomNumber(roomNumber string, roomID string) error {
	key := fmt.Sprintf("room:%s", roomNumber)
	_, err := r.kvs.Set(context.Background(), key, roomID, 0).Result()
	return err
}

func (r *sessionRepository) DeleteRoomNumber(roomNumber string) error {
	key := fmt.Sprintf("room:%s", roomNumber)
	_, err := r.kvs.Del(context.Background(), key).Result()
	return err
}
