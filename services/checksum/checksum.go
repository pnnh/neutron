package checksum

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalcFileHash(path string) (hash string, retErr error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			retErr = fmt.Errorf("error closing file %s: %v; original error: %w", path, closeErr, retErr)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
