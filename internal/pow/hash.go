package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
)

const (
	nonceSize  = 16
	tokenSize  = 32
	targetBits = 64
)

var (
	ErrInvalidComplexity = errors.New("invalid complexity")
	ErrInvalidChallenge  = errors.New("invalid challenge")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrUnverified        = errors.New("unverified")
)

type Pow struct {
	complexity uint64
}

func New(complexity uint64) (*Pow, error) {
	if complexity < 1 || complexity > targetBits {
		return nil, ErrInvalidComplexity
	}
	return &Pow{complexity: complexity}, nil
}

func (p *Pow) Challenge() []byte {
	return newToken(p.complexity)
}

func newToken(targetBits uint64) []byte {
	buf := make([]byte, tokenSize)
	target := uint64(1) << (64 - targetBits)

	binary.BigEndian.PutUint64(buf[:8], target)
	_, _ = rand.Read(buf[8:])

	return buf
}

func (p *Pow) Verify(challenge, solution []byte) error {
	if len(challenge) != tokenSize {
		return ErrInvalidChallenge
	}

	if len(solution) != nonceSize {
		return ErrInvalidSolution
	}

	if !verify(challenge, solution) {
		return ErrUnverified
	}

	return nil
}

func verify(token, nonce []byte) bool {
	h := hash(token, nonce)
	return bytes.Compare(h, token) < 0
}

func hash(token, nonce []byte) []byte {
	h := sha256.New()
	h.Write(token)
	h.Write(nonce)
	return h.Sum(nil)
}

func (p *Pow) Solve(challenge []byte) []byte {
	if len(challenge) != tokenSize {
		return nil
	}

	return solve(challenge)
}

func solve(token []byte) []byte {
	nonce := make([]byte, nonceSize)

	for i := uint64(0); i < math.MaxUint64; i++ {
		binary.BigEndian.PutUint64(nonce, i)
		if verify(token, nonce) {
			return nonce
		}
	}

	return nil
}
