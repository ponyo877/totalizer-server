package session

type Reader interface {
}

type Writer interface {
	Incriment(string) (int, error)
}

type Repository interface {
	Writer
	Reader
}

type UseCase interface {
	Incriment(string) (int, error)
}
