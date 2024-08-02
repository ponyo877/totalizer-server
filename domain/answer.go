package domain

import (
	"encoding/json"
)

type Type int

const (
	TypeEnter Type = iota
	TypeQuestion
	TypeReady
	TypeResult
)

type Answer interface {
	String() (string, error)
}

func answerToString(a Answer) (string, error) {
	json, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

type BaseAnswer struct {
	Type Type   `json:"type"`
	From string `json:"from"`
}

type EnterAnswer struct {
	BaseAnswer
	EnterCount int `json:"enter_count"`
}

type AskAnswer struct {
	BaseAnswer
	QuestionID      string `json:"question_id"`
	QuestionContent string `json:"question_content"`
}

type ReadyAnswer struct {
	BaseAnswer
}

type ReleaseAnswer struct {
	BaseAnswer
	YesCount   int `json:"yes_count"`
	EnterCount int `json:"enter_count"`
}

type OpenAnswer struct {
	BaseAnswer
	RoomID     string `json:"room_id"`
	RoomNumber string `json:"room_number"`
}

type StatsAnswer struct {
	BaseAnswer
	RoomID          string `json:"room_id"`
	EnterCount      int    `json:"enter_count,omitempty"`
	QuestionID      string `json:"question_id,omitempty"`
	QuestionContent string `json:"question_content,omitempty"`
	YesCount        *int   `json:"yes_count,omitempty"`
}

func NewEnterAnswer(enterCount int) (*EnterAnswer, error) {
	base := BaseAnswer{
		Type: TypeEnter,
		From: "user",
	}
	return &EnterAnswer{base, enterCount}, nil
}

func NewAskAnswer(quesionID, questionContent string) (*AskAnswer, error) {
	base := BaseAnswer{
		Type: TypeQuestion,
		From: "user",
	}
	return &AskAnswer{base, quesionID, questionContent}, nil
}

func NewReadyAnswer() (*ReadyAnswer, error) {
	base := BaseAnswer{
		Type: TypeReady,
		From: "user",
	}
	return &ReadyAnswer{base}, nil
}

func NewResultAnswer(yesCount, enterCount int) (*ReleaseAnswer, error) {
	base := BaseAnswer{
		Type: TypeResult,
		From: "user",
	}
	return &ReleaseAnswer{base, yesCount, enterCount}, nil
}

func NewOpenAnswer(roomID, roomNumber string) (*OpenAnswer, error) {
	base := BaseAnswer{
		Type: TypeEnter,
		From: "system",
	}
	return &OpenAnswer{
		BaseAnswer: base,
		RoomID:     roomID,
		RoomNumber: roomNumber,
	}, nil
}

func NewStatsAnswer(roomID string, enterCount int, questionID string, questionContent string, yesCount *int) (*StatsAnswer, error) {
	base := BaseAnswer{
		Type: TypeEnter,
		From: "system",
	}
	return &StatsAnswer{
		BaseAnswer:      base,
		RoomID:          roomID,
		EnterCount:      enterCount,
		QuestionID:      questionID,
		QuestionContent: questionContent,
		YesCount:        yesCount,
	}, nil
}

func (a *EnterAnswer) String() (string, error) {
	return answerToString(a)
}

func (a *AskAnswer) String() (string, error) {
	return answerToString(a)
}

func (a *ReadyAnswer) String() (string, error) {
	return answerToString(a)
}

func (a *ReleaseAnswer) String() (string, error) {
	return answerToString(a)
}
