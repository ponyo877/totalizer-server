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

func (s *Status) Stats(ec int, qid, qcont string, vr int) *Stats {
	switch s.status {
	// 人数
	case StatusOpen:
		return NewStats(ec, "", "", 0)
	// 人数, 質問
	case StatusQuestion:
		return NewStats(ec, qid, qcont, 0)
	// 人数, 質問
	case StatusReady:
		return NewStats(ec, qid, qcont, 0)
	// 人数, 質問, 投票結果
	case StatusResult:
		return NewStats(ec, qid, qcont, vr)
	default:
		return NewStats(0, "", "", 0)
	}
}
