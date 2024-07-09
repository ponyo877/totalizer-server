package session

import "github.com/ponyo877/totalizer-server/domain"

type Service struct {
	repository Repository
}

func NewService(r Repository) UseCase {
	return &Service{
		repository: r,
	}
}

func (s *Service) Incriment(key string) (int, error) {
	return s.repository.Incriment(key)
}

func (s *Service) ListQuestion() (*domain.Question, error) {
	return s.repository.ListQuestion()
}

func (s *Service) OpenRoom(roomID string) error {
	return nil
}

func (s *Service) Subscribe(roomID string) error {
	return nil
}

func (s *Service) IncEnterCount(roomID string) (int, error) {
	return 0, nil
}

func (s *Service) RegisterQuestion(roomID string, question string) error {
	return nil
}

func (s *Service) PublishQuestion(roomID string, question string) error {
	return nil
}

func (s *Service) CountVote(roomID string) (int, error) {
	return 0, nil
}

func (s *Service) FetchEnterCount(roomID string) (int, error) {
	return 0, nil
}

func (s *Service) VoteYes(roomID string) error {
	return nil
}

func (s *Service) PublishReady(roomID string) error {
	return nil
}

func (s *Service) PublishResult(roomID string) error {
	return nil
}

func (s *Service) SaveResult(roomID string, questionID string) error {
	return nil
}
