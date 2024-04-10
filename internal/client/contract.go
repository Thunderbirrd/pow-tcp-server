package client

type pow interface {
	solver
}

type solver interface {
	Solve(challenge []byte) []byte
}
