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
