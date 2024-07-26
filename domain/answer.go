package domain

import (
	"encoding/json"
	"errors"
)

type AnswerType int

const (
	AnswerTypeEnter AnswerType = iota
	AnswerTypeQuestion
	AnswerTypeReady
	AnswerTypeResult
)

type Answer struct {
	answerType AnswerType
	value      interface{}
}

func NewAnswer(answerType AnswerType, value interface{}) (*Answer, error) {
	ans := &Answer{answerType, value}
	if err := ans.validation(); err != nil {
		return nil, err
	}
	return ans, nil
}

func (a *Answer) validation() error {
	switch a.answerType {
	case AnswerTypeEnter:
		if a.value != nil {
			return errors.New("Content is not nil")
		}
	case AnswerTypeQuestion:
		if _, ok := a.value.(string); !ok {
			return errors.New("Content is not string")
		}
	case AnswerTypeReady:
		if a.value != nil {
			return errors.New("Content is not nil")
		}
	case AnswerTypeResult:
		if _, ok := a.value.(int); !ok {
			return errors.New("Content is not int")
		}
	default:
		return errors.New("AnswerType is invalid")
	}
	return nil
}

func (a *Answer) Type() AnswerType {
	return a.answerType
}

func (a *Answer) Quesion() string {
	return a.value.(string)
}

func (a *Answer) Result() int {
	return a.value.(int)
}

func (a *Answer) String() string {
	js := struct {
		Type  AnswerType  `json:"type"`
		Value interface{} `json:"value,omitempty"`
	}{
		Type:  a.answerType,
		Value: a.value,
	}
	json, _ := json.Marshal(js)
	return string(json)
}
