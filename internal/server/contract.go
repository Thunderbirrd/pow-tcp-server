package server

type pow interface {
	verifier
}

type verifier interface {
	Challenge() []byte
	Verify(challenge, solution []byte) error
}

type repo interface {
	GetLine() (string, error)
}
