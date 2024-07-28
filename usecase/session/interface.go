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
	PublishResult(string, string) error
	UpdateQuestionVote(string) error
	PublishEnter(string) error
	StoreRoomStatus(roomID string, status domain.RoomStatus) error
	GetRoomStatus(string) (*domain.Status, error)
	GetLatestQuestion(string) (*domain.Question, error)
}

type UseCase interface {
	Enter(string) (*chan string, error)
	Ask(string, string) error
	Vote(string, string, string) error
	Release(string, string) error
	FetchStats(roomID string) (*domain.Stats, error)
}
