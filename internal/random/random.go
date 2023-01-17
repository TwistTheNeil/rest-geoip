package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateKey generates a crypto secure key
// Thanks rclone!
// https://github.com/rclone/rclone/blob/master/lib/random/random.go
func GenerateKey(bits int) (string, error) {
	bytes := bits / 8
	if bits%8 != 0 {
		bytes++
	}
	var pw = make([]byte, bytes)
	n, err := rand.Read(pw)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if n != bytes {
		return "", fmt.Errorf("%w", fmt.Errorf("who knows"))
	}
	return base64.RawURLEncoding.EncodeToString(pw), nil
}
