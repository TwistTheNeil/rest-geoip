package fileops

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"rest-geoip/lib/customerrors"
)

func md5Hash(file string) ([]byte, error) {
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

// VerifyMD5HashFromFile hashes a file and verifies it against a sum
// contained within a .md5 file
func VerifyMD5HashFromFile(file, md5sumFile string) error {
	actual, err := md5Hash(file)
	if err != nil {
		return err
	}

	cleanMD5SumFile := filepath.Clean(md5sumFile)

	// We know exactly where this file and path is
	// #nosec G304
	expected, err := ioutil.ReadFile(cleanMD5SumFile)
	if err != nil {
		return err
	}

	if fmt.Sprintf("%x", actual) != fmt.Sprintf("%s", expected) {
		return customerrors.ErrHashChecksum
	}

	return nil
}
