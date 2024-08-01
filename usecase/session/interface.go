package session

import "github.com/ponyo877/totalizer-server/domain"

type Repository interface {
	SubscribeRoom(string) *chan string
	IncrimentEnterCount(string) (int, error)
	CreateQuestion(*domain.Question) error
	PublishQuestion(*domain.Question) error
	GetVoteCount(string) (int, error)
	GetAnswerCount(string, string) (int, error)
	GetEnterCount(string) (int, error)
	IncrimentVoteCount(string, string) (int, error)
	PublishReady(string) error
	PublishResult(string, int, int) error
	UpdateQuestionVote(string) error
	PublishEnter(string, int) error
	StoreRoomStatus(string, domain.RoomStatus) error
	GetRoomStatus(string) (*domain.Status, error)
	GetLatestQuestion(string) (*domain.Question, error)
	GetRoomIDByRoomNumber(string) (string, bool, error)
	SetRoomNumber(string, string) error
	DeleteRoomNumber(string) error
}

type UseCase interface {
	Open(string) (string, error)
	Enter(string) (*chan string, error)
	Ask(string, string) error
	Vote(string, string, string) error
	Release(string, string) error
	FetchStats(string) (*domain.Stats, error)
	FetchRoomID(string) (string, error)
}
