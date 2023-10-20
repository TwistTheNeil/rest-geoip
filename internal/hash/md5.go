package hash

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
)

func MD5Hash(file string) ([]byte, error) {
	filePath := filepath.Clean(file)

	// We know exactly where this file and path is
	// #nosec G304
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// #nosec G401
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), f.Close()
}
