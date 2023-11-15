package util

import "crypto/sha256"

// HashSHA256 hashes the given string with SHA256
func HashSHA256(s string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(s))

	return hasher.Sum(nil)
}
