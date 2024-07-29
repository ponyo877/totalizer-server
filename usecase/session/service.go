package session

import (
	"errors"
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
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return nil, err
	}
	if status == nil {
		if err := s.repository.StoreRoomStatus(roomID, domain.StatusOpen); err != nil {
			return nil, err
		}
	}
	ch := s.repository.SubscribeRoom(roomID)
	enterCount, err := s.repository.IncrimentEnterCount(roomID)
	if err != nil {
		return nil, err
	}
	if err := s.repository.PublishEnter(roomID, enterCount); err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *Service) Ask(roomID string, question string) error {
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return err
	}
	if !status.IsOpen() {
		return errors.New("room status is invalid")
	}
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
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return err
	}
	if !status.IsQuestion() {
		return errors.New("room status is invalid")
	}
	if _, err := s.repository.IncrimentVoteCount(questionID, answer); err != nil {
		return err
	}
	count, err := s.repository.GetVoteCount(questionID)
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
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return err
	}
	if !status.IsReady() {
		return errors.New("room status is invalid")
	}
	if err := s.repository.UpdateQuestionVote(questionID); err != nil {
		return err
	}
	if err := s.repository.StoreRoomStatus(roomID, domain.StatusResult); err != nil {
		return err
	}
	yesCount, err := s.repository.GetAnswerCount(questionID, "yes")
	if err != nil {
		return err
	}
	enterCount, err := s.repository.GetEnterCount(roomID)
	if err != nil {
		return err
	}
	return s.repository.PublishResult(roomID, yesCount, enterCount)
}

func (s *Service) FetchStats(roomID string) (*domain.Stats, error) {
	status, err := s.repository.GetRoomStatus(roomID)
	if err != nil {
		return nil, err
	}
	var enterCount, yesCount int
	var questionID, questionContent string
	enterCount, err = s.repository.GetEnterCount(roomID)
	if err != nil {
		return nil, err
	}
	question, err := s.repository.GetLatestQuestion(roomID)
	if err != nil {
		return nil, err
	}
	if question != nil {
		questionID = question.ID()
		questionContent = question.Content()
		yesCount, err = s.repository.GetAnswerCount(question.ID(), "yes")
		if err != nil {
			return nil, err
		}
	}
	return status.Stats(enterCount, questionID, questionContent, yesCount), nil
}
