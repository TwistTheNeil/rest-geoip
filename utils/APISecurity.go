package utils

import (
	"crypto/rand"
	"encoding/base64"
	"rest-geoip/customerrors"
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
		return "", customerrors.ErrGeneratePassword
	}
	if n != bytes {
		return "", customerrors.ErrGeneratePassword
	}
	return base64.RawURLEncoding.EncodeToString(pw), nil
}
