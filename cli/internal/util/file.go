package util

import (
	"github.com/charmbracelet/log"
	"os"
)

// ReadFile read the contents of a file at path p and return string contents
func ReadFile(p string) (string, error) {
	dat, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

// WriteFile write contents to file
func WriteFile(data []byte, path string) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error("failed to close file", "err", err)
		}
	}(f)
	_, err = f.Write(data)
	return err
}
