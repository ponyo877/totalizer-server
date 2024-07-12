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
	return s.repository.OpenRoom(roomID)
}

func (s *Service) Subscribe(roomID string) error {
	return s.repository.SubscribeRoom(roomID)
}

func (s *Service) IncEnterCount(roomID string) (int, error) {
	return s.repository.IncrimentEnterCount(roomID)
}

func (s *Service) RegisterQuestion(roomID string, question string) error {
	return s.repository.CreateQuestion(roomID, question)
}

func (s *Service) PublishQuestion(roomID string, question string) error {
	return s.repository.PublishQuestion(roomID, question)
}

func (s *Service) CountVote(roomID string) (int, error) {
	return s.repository.GetVoteCount(roomID)
}

func (s *Service) FetchEnterCount(roomID string) (int, error) {
	return s.repository.GetEnterCount(roomID)
}

func (s *Service) VoteYes(roomID string) error {
	return s.repository.IncrimentVoteCount(roomID, "YES")
}

func (s *Service) PublishReady(roomID string) error {
	return s.repository.PublishReady(roomID)
}

func (s *Service) PublishResult(roomID string, question string) error {
	return s.repository.PublishResult(roomID, question)
}

func (s *Service) SaveResult(roomID string, questionID string) error {
	return s.repository.UpdateQuestionVote(roomID, questionID)
}
