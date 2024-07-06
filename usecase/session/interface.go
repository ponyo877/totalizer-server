package session

import "github.com/ponyo877/totalizer-server/domain"

type Repository interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
}

type UseCase interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
}
