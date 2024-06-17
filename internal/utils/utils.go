package utils

import (
	"crypto/rand"
	"math/big"
)

func RandomInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		Logger.Fatalf("Failed to generate random number: %v", err)
	}
	return int(nBig.Int64())
}
