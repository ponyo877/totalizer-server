package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/ponyo877/totalizer-server/domain"
)

type Service struct {
	repository Repository
}

func NewService(r Repository) UseCase {
	return &Service{
		repository: r,
	}
}

func (s *Service) Enter(roomID string) (*chan string, error) {
	ch := s.repository.SubscribeRoom(roomID)
	if err := s.repository.PublishEnter(roomID); err != nil {
		return nil, err
	}
	if _, err := s.repository.IncrimentEnterCount(roomID); err != nil {
		return nil, err
	}
	if err := s.repository.StoreRoomStatus(roomID, domain.StatusOpen); err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *Service) Ask(roomID string, question string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	q := domain.NewQuestion(id.String(), roomID, question, 0, time.Now())
	if err := s.repository.CreateQuestion(q); err != nil {
		return err
	}
	if err := s.repository.PublishQuestion(q); err != nil {
		return err
	}
	if err := s.repository.StoreRoomStatus(roomID, domain.StatusQuestion); err != nil {
		return err
	}
	return nil
}

func (s *Service) Vote(roomID string, questionID string, answer string) error {
	if answer == "yes" {
		if _, err := s.repository.IncrimentVoteCount(roomID, "YES"); err != nil {
			return err
		}
	}
	count, err := s.repository.GetVoteCount(roomID)
	if err != nil {
		return err
	}
	total, err := s.repository.GetEnterCount(roomID)
	if err != nil {
		return err
	}
	if count == total {
		if err := s.repository.PublishReady(roomID); err != nil {
			return err
		}
		if err := s.repository.StoreRoomStatus(roomID, domain.StatusReady); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Release(roomID string, questionID string) error {
	if err := s.repository.UpdateQuestionVote(questionID); err != nil {
		return err
	}
	if err := s.repository.StoreRoomStatus(roomID, domain.StatusResult); err != nil {
		return err
	}
	return s.repository.PublishResult(roomID, questionID)
}

func (s *Service) Sync(roomID string) (*domain.Status, error) {
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return nil, err
	}
	// 人数
	if status.IsOpen() {
		return nil, nil
	}
	// 人数, 質問, 回答人数
	if status.IsQuestion() {
		return nil, nil
	}
	// 人数, 質問
	if status.IsReady() {
		return nil, nil
	}
	// 人数, 質問, 回答結果
	if status.IsResult() {
		return nil, nil
	}
	return nil, nil
}
