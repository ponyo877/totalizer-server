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

func (s *Socket) Open(roomID string) error {
	mu[roomID] = &sync.Mutex{}
	if err := s.service.Open(roomID); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Enter(roomID string) error {
	mu[roomID].Lock()
	defer mu[roomID].Unlock()
	if err := s.service.Enter(roomID); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Ask(roomID string, question string) error {
	if err := s.service.Ask(roomID, question); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Vote(roomID string, questionID string, answer string) error {
	mu[roomID].Lock()
	defer mu[roomID].Unlock()
	if err := s.service.Vote(roomID, questionID, answer); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Release(roomID string, questionID string) error {
	if err := s.service.Release(roomID, questionID); err != nil {
		return err
	}
	return s.send("ok")
}
