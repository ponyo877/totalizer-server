package session

import "github.com/ponyo877/totalizer-server/domain"

type Service struct {
	repository Repository
}

func NewService(r Repository) *Service {
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
