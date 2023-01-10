package fileops

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

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
