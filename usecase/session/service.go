package session

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
