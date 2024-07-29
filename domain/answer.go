package domain

import (
	"encoding/json"
)

type AnswerType int

const (
	AnswerTypeEnter AnswerType = iota
	AnswerTypeQuestion
	AnswerTypeReady
	AnswerTypeResult
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

type EnterAnswer struct {
	AnswerType AnswerType `json:"type"`
	EnterCount int        `json:"enter_count"`
}

type AskAnswer struct {
	AnswerType      AnswerType `json:"type"`
	QuestionID      string     `json:"question_id"`
	QuestionContent string     `json:"question_content"`
}

type ReadyAnswer struct {
	AnswerType AnswerType `json:"type"`
}

type ReleaseAnswer struct {
	AnswerType AnswerType `json:"type"`
	YesCount   int        `json:"yes_count"`
	EnterCount int        `json:"enter_count"`
}

func NewEnterAnswer(enterCount int) (*EnterAnswer, error) {
	return &EnterAnswer{AnswerTypeEnter, enterCount}, nil
}

func NewAskAnswer(quesionID, questionContent string) (*AskAnswer, error) {
	return &AskAnswer{AnswerTypeQuestion, quesionID, questionContent}, nil
}

func NewReadyAnswer() (*ReadyAnswer, error) {
	return &ReadyAnswer{AnswerTypeReady}, nil
}

func NewResultAnswer(yesCount, enterCount int) (*ReleaseAnswer, error) {
	return &ReleaseAnswer{AnswerTypeResult, yesCount, enterCount}, nil
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
