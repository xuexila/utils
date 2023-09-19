package utils

import "crypto/sha256"

// sha256 hash
func Sha256(msg string) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(msg))
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}
