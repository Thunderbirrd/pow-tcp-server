package repository

import (
	"bufio"
	"bytes"
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed wow.txt
var lines []byte

type Repo struct {
	lines []string
}

func New() *Repo {
	repo := new(Repo)
	reader := bytes.NewReader(lines)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if q := strings.TrimSpace(scanner.Text()); q != "" {
			repo.lines = append(repo.lines, q)
		}
	}

	return repo
}

func (q *Repo) GetLine() (string, error) {
	return q.lines[rand.Intn(len(q.lines))], nil
}
