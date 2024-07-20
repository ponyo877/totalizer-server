package domain

type Status struct {
	// OPEN, QUESION, READY, RESULT
	status string
}

func NewStatus(status string) *Status {
	return &Status{status}
}

func (s *Status) IsOpen() bool {
	return s.status == "OPEN"
}

func (s *Status) IsQuestion() bool {
	return s.status == "QUESTION"
}

func (s *Status) IsReady() bool {
	return s.status == "READY"
}

func (s *Status) IsResult() bool {
	return s.status == "RESULT"
}
