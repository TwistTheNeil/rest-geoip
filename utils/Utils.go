package utils

import (
	"archive/tar"
	"compress/gzip"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"

	// We don't have a choice here
	// #nosec G501
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"rest-geoip/customerrors"
	"sync"

	"github.com/markbates/pkger"
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

// FindFile returns a path to a file matching regex under root
// Returns
//  string: Full path
//  string: File name
//  error : Error
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

// ExtractTarGz extracts a gzipped stream to dest
func ExtractTarGz(r io.Reader, dest string) error {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("Stream requires gzip-compressed body: %v", err)
	}

	tr := tar.NewReader(zr)

	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Tar error: %v", err)
		}

		switch f.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dest+f.Name, 0750); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %v", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(dest + f.Name)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %v", err)
			}
			// For our purposes, we don't expect any files larger than 100MiB
			limited := &io.LimitedReader{R: tr, N: 100 << 20}
			if _, err := io.Copy(outFile, limited); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %v", err)
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		default:
			return fmt.Errorf(
				"ExtractTarGz: %s has uknown type: %v",
				f.Name,
				f.Typeflag)
		}
	}

	return nil
}

// ParseTemplates parses all templates
func ParseTemplates(root string) (*template.Template, error) {
	t := template.New("")

	// Since Walk receives a dynamic value, pkger won't be able to find the
	// actual directory to package from the next line, which is why we used
	// pkger.Include() in main().
	err := pkger.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		// We read from pkger's fs here so the template can be parsed
		contents, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = t.Parse(string(contents))
		if err != nil {
			return err
		}
		return nil
	})

	return t, err
}
