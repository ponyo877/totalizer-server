package domain

import "time"

type Question struct {
	id       string
	roomID   string
	content  string
	vote     int
	createAt time.Time
}

func NewQuestion(id, roomID, content string, vote int, createAt time.Time) *Question {
	return &Question{id, roomID, content, vote, createAt}
}

func (q *Question) ID() string {
	return q.id
}

func (q *Question) RoomID() string {
	return q.roomID
}

func (q *Question) Content() string {
	return q.content
}

func (q *Question) Vote() int {
	return q.vote
}

func (q *Question) CreatedAt() time.Time {
	return q.createAt
}
