package domain

import "time"

type Question struct {
	id        string
	roomID    string
	content   string
	voteCount int
	yesCount  int
	createAt  time.Time
}

func NewQuestion(id, roomID, content string, voteCount, yesCount int, createAt time.Time) *Question {
	return &Question{id, roomID, content, voteCount, yesCount, createAt}
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

func (q *Question) VoteCount() int {
	return q.voteCount
}

func (q *Question) YesCount() int {
	return q.yesCount
}

func (q *Question) CreatedAt() time.Time {
	return q.createAt
}
