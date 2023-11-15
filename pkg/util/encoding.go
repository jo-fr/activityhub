package util

import "encoding/base64"

// DecodeBase64 decodes the given base64 string
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
