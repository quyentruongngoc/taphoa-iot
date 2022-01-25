package internal

import (
	"crypto/sha256"
	"fmt"
)

func SHA256(data string) string {
	b := []byte(data)
	hash := sha256.Sum256(b)
	ret := fmt.Sprintf("%x", hash)

	return ret
}
