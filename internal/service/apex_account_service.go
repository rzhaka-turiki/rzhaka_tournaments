package service

import (
	"crypto/sha512"
	"encoding/hex"
)

func calculateNIDHash(uid string) string {
	hash := sha512.Sum512([]byte(uid))
	return hex.EncodeToString(hash[:])
}
