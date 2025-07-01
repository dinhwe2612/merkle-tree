package tree

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func HashConcat(a, b string) string {
	h := sha256.Sum256([]byte(a + b))
	return hex.EncodeToString(h[:])
}

type ProofStep struct {
	key    string
	isLeft bool
}

var TestCases = []struct {
	name       string
	numLeaves  int
	numQueries int
}{
	{"Small", 1000, 1000},
	{"Medium", 100000, 100000},
	{"Large", 5000000, 5000000},
}
