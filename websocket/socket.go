package socket

import (
	"github.com/google/uuid"
	"github.com/ponyo877/totalizer-server/domain"
	"github.com/ponyo877/totalizer-server/usecase/session"
	"golang.org/x/net/websocket"
)

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
	roomNumber, err := s.service.Open(roomID)
	if err != nil {
		return err
	}
	ch, err := s.service.Enter(roomID)
	if err != nil {
		return err
	}
	go s.recieve(ch)
	answer, _ := domain.NewOpenAnswer(roomID, roomNumber)
	return s.sendJSON(answer)
}

func (s *Socket) Enter(roomNumber string) error {
	roomID, err := s.service.FetchRoomID(roomNumber)
	if err != nil {
		s.sendJSON(domain.NewFailAnswer())
		return err
	}
	ch, err := s.service.Enter(roomID)
	if err != nil {
		s.sendJSON(domain.NewFailAnswer())
		return err
	}
	go s.recieve(ch)
	stats, err := s.service.FetchStats(roomID)
	if err != nil {
		s.sendJSON(domain.NewFailAnswer())
		return err
	}
	answer, _ := domain.NewStatsAnswer(
		roomID,
		stats.EnterCount(),
		stats.QuestionID(),
		stats.QuestionContent(),
		stats.YesCount(),
	)
	return s.sendJSON(answer)
}

func (s *Socket) Ask(roomID string, question string) error {
	if err := s.service.Ask(roomID, question); err != nil {
		return err
	}
	return nil
}

func (s *Socket) Vote(roomID string, questionID string, answer string) error {
	if err := s.service.Vote(roomID, questionID, answer); err != nil {
		return err
	}
	return nil
}

func (s *Socket) Release(roomID string, questionID string) error {
	if err := s.service.Release(roomID, questionID); err != nil {
		return err
	}
	return nil
}
