package session

import "github.com/ponyo877/totalizer-server/domain"

type Repository interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
	OpenRoom(string) error
	SubscribeRoom(string) error
	IncrimentEnterCount(string) (int, error)
	CreateQuestion(string, string) error
	PublishQuestion(string, string) error
	GetVoteCount(string) (int, error)
	GetEnterCount(string) (int, error)
	IncrimentVoteCount(string, string) error
	PublishReady(string) error
	PublishResult(string, string) error
	UpdateQuestionVote(string, string) error
}

type UseCase interface {
	Incriment(string) (int, error)
	ListQuestion() (*domain.Question, error)
	OpenRoom(string) error
	Subscribe(string) error
	IncEnterCount(string) (int, error)
	RegisterQuestion(string, string) error
	PublishQuestion(string, string) error
	CountVote(string) (int, error)
	FetchEnterCount(string) (int, error)
	VoteYes(string) error
	PublishReady(string) error
	PublishResult(string, string) error
	SaveResult(string, string) error
}
