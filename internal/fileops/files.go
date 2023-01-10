package fileops

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"rest-geoip/lib/customerrors"
	"sync"
)

// Download a file
func Download(url, dest string, wg *sync.WaitGroup, errChannel chan<- error) {
	defer wg.Done()

	// We know exactly how this url is constructed
	// #nosec G107
	resp, err := http.Get(url)
	if err != nil {
		errChannel <- customerrors.ErrDownloadFile
		return
	}
	if resp.StatusCode%200 > 99 {
		errChannel <- fmt.Errorf("Download error: %s", resp.Status)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		errChannel <- customerrors.ErrCreateFile
		return
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		errChannel <- customerrors.ErrCreateFile
		return
	}

	if err = out.Close(); err != nil {
		errChannel <- err
	}
}

// FindFile returns a path to a file matching regex under root
// Returns
//
//	string: Full path
//	string: File name
//	error : Error
func FindFile(root, r string) (string, string, error) {
	regex, err := regexp.Compile(r)
	if err != nil {
		return "", "", err
	}

	var foundPath string
	var foundName string

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if regex.MatchString(info.Name()) {
			foundPath = path
			foundName = info.Name()
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}
	return foundPath, foundName, nil
}

// MoveFile moves a file
func MoveFile(source, dest string) error {
	// #nosec G304
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, input, 0600)
	if err != nil {
		return err
	}

	return nil
}
