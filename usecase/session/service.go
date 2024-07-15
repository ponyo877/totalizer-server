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

func (s *Service) Open(roomID string) error {
	if err := s.repository.OpenRoom(roomID); err != nil {
		return err
	}
	if err := s.repository.SubscribeRoom(roomID); err != nil {
		return err
	}
	return nil
}

func (s *Service) Enter(roomID string) error {
	if err := s.repository.SubscribeRoom(roomID); err != nil {
		return err
	}
	if _, err := s.repository.IncrimentEnterCount(roomID); err != nil {
		return err
	}
	return nil
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
	}
	return nil
}

func (s *Service) Release(roomID string, questionID string) error {
	if err := s.repository.UpdateQuestionVote(questionID); err != nil {
		return err
	}
	return s.repository.PublishResult(roomID, questionID)
}
