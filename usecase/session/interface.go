package session

import "github.com/ponyo877/totalizer-server/domain"

type Repository interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
}

type UseCase interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
	OpenRoom(string) error
	Subscribe(string) error
	IncEnterCount(string) (int, error)
	RegisterQuestion(string, string) error
	PublishQuestion(string, string) error
	VoteAny(string) (int, error)
	FetchEnterCount(string) (int, error)
	VoteYes(string) error
	PublishReady(string) error
	PublishResult(string) error
	SaveResult(string, string) error
}
