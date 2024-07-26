package domain

type RoomStatus int

const (
	StatusOpen RoomStatus = iota
	StatusQuestion
	StatusReady
	StatusResult
)

type Status struct {
	status RoomStatus
}

func NewStatus(status RoomStatus) *Status {
	return &Status{status}
}

func (s *Status) IsOpen() bool {
	return s.status == StatusOpen
}

func (s *Status) IsQuestion() bool {
	return s.status == StatusQuestion
}

func (s *Status) IsReady() bool {
	return s.status == StatusReady
}

func (s *Status) IsResult() bool {
	return s.status == StatusResult
}
