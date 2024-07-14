package session

import "github.com/ponyo877/totalizer-server/domain"

type Repository interface {
	OpenRoom(string) error
	SubscribeRoom(string) error
	IncrimentEnterCount(string) (int, error)
	CreateQuestion(*domain.Question) error
	PublishQuestion(*domain.Question) error
	GetVoteCount(string) (int, error)
	GetEnterCount(string) (int, error)
	IncrimentVoteCount(string, string) (int, error)
	PublishReady(string) error
	PublishResult(string, string) error
	UpdateQuestionVote(string) error
}

type UseCase interface {
	Open(string) error
	Enter(string) error
	Ask(string, string) error
	Vote(string, string, string) error
	Release(string, string) error
}
