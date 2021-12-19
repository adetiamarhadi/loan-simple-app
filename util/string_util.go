package util

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomAlphaNumeric return random string alphanumeric
func RandomAlphaNumeric(length int) string {
	return stringWithCharset(length, charset)
}