package domain

type Stats struct {
	enterCount      int
	questionID      string
	questionContent string
	yesCount        *int
}

func NewStats(enterCount int, questionID, questionContent string, yesCount *int) *Stats {
	return &Stats{
		enterCount:      enterCount,
		questionID:      questionID,
		questionContent: questionContent,
		yesCount:        yesCount,
	}
}

func (s *Stats) EnterCount() int {
	return s.enterCount
}

func (s *Stats) QuestionID() string {
	return s.questionID
}

func (s *Stats) QuestionContent() string {
	return s.questionContent
}

func (s *Stats) YesCount() *int {
	return s.yesCount
}
