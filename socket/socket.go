package socket

import (
	"sync"

	"github.com/ponyo877/totalizer-server/usecase/session"
	"golang.org/x/net/websocket"
)

var mu map[string]*sync.Mutex

type Socket struct {
	ws      *websocket.Conn
	service session.UseCase
}

func NewSocket(ws *websocket.Conn, service session.UseCase) *Socket {
	return &Socket{ws, service}
}

func (s *Socket) send(msg interface{}) error {
	return websocket.JSON.Send(s.ws, msg)
}

func (s *Socket) receive(msg interface{}) error {
	return websocket.JSON.Receive(s.ws, msg)
}

func (s *Socket) Open(roomID string) error {
	mu[roomID] = &sync.Mutex{}
	if err := s.service.OpenRoom(roomID); err != nil {
		return err
	}
	if err := s.service.Subscribe(roomID); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Enter(roomID string) error {
	if err := s.service.Subscribe(roomID); err != nil {
		return err
	}
	if _, err := s.service.IncEnterCount(roomID); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Ask(roomID string, question string) error {
	if err := s.service.RegisterQuestion(roomID, question); err != nil {
		return err
	}
	if err := s.service.PublishQuestion(roomID, question); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Vote(roomID string, questionID string, answer string) error {
	mu[roomID].Lock()
	defer mu[roomID].Unlock()
	if answer == "yes" {
		if err := s.service.VoteYes(questionID); err != nil {
			return err
		}
	}
	count, err := s.service.CountVote(questionID)
	if err != nil {
		return err
	}
	total, err := s.service.FetchEnterCount(roomID)
	if err != nil {
		return err
	}
	if count == total {
		if err := s.service.PublishReady(roomID); err != nil {
			return err
		}
	}
	return s.send("ok")
}

func (s *Socket) Release(roomID string, questionID string) error {
	if err := s.service.SaveResult(roomID, questionID); err != nil {
		return err
	}
	return s.service.PublishResult(roomID)
}
