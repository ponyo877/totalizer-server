package socket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/ponyo877/totalizer-server/usecase/session"
	"golang.org/x/net/websocket"
)

var mu = map[string]*sync.Mutex{}

type Socket struct {
	ws      *websocket.Conn
	service session.UseCase
}

func NewSocket(ws *websocket.Conn, service session.UseCase) *Socket {
	return &Socket{ws, service}
}

func (s *Socket) send(msg interface{}) error {
	return websocket.Message.Send(s.ws, msg)
}

func (s *Socket) sendJSON(msg interface{}) error {
	return websocket.JSON.Send(s.ws, msg)
}

func (s *Socket) recieve(ch *chan string) {
	for msg := range *ch {
		s.send(msg)
	}
}

func (s *Socket) Open() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	roomID := uuid.String()
	mu[roomID] = &sync.Mutex{}
	ch, err := s.service.Enter(roomID)
	if err != nil {
		return err
	}
	go s.recieve(ch)
	return s.send(roomID)
}

func (s *Socket) Enter(roomID string) error {
	if _, ok := mu[roomID]; !ok {
		return s.send("room not found")
	}
	mu[roomID].Lock()
	defer mu[roomID].Unlock()
	ch, err := s.service.Enter(roomID)
	if err != nil {
		return err
	}
	go s.recieve(ch)
	return s.send("ok")
}

func (s *Socket) Ask(roomID string, question string) error {
	if err := s.service.Ask(roomID, question); err != nil {
		return err
	}
	return s.send("ok")
}

func (s *Socket) Vote(roomID string, questionID string, answer string) error {
	if _, ok := mu[roomID]; !ok {
		return s.send("room not found")
	}
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
