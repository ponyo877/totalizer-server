package domain

import (
	"encoding/json"
	"errors"
)

type Answer struct {
	// ENTER, QUESION, READY, RESULT
	answerType string
	value      interface{}
}

func NewAnswer(answerType string, value interface{}) (*Answer, error) {
	ans := &Answer{answerType, value}
	if err := ans.validation(); err != nil {
		return nil, err
	}
	return ans, nil
}

func (a *Answer) validation() error {
	switch a.answerType {
	case "ENTER":
		if a.value != nil {
			return errors.New("Content is not nil")
		}
	case "QUESTION":
		if _, ok := a.value.(string); !ok {
			return errors.New("Content is not string")
		}
	case "READY":
		if a.value != nil {
			return errors.New("Content is not nil")
		}
	case "RESULT":
		if _, ok := a.value.(int); !ok {
			return errors.New("Content is not int")
		}
	default:
		return errors.New("AnswerType is invalid")
	}
	return nil
}

func (a *Answer) Type() string {
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
		Type  string      `json:"type"`
		Value interface{} `json:"value,omitempty"`
	}{
		Type:  a.answerType,
		Value: a.value,
	}
	json, _ := json.Marshal(js)
	return string(json)
}
