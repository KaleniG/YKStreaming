package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
)

func GenerateRandomToken(length uint32) (string, error) {
	if length == 0 {
		return "", nil
	}

	nBytes := (length + 1) / 2

	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	s := hex.EncodeToString(b)

	return s[:length], nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FindFilesContaining(dir, substr string) ([]string, error) {
	var result []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.Contains(entry.Name(), substr) {
			result = append(result, entry.Name())
		}
	}

	return result, nil
}

func GetEnvDir(name string) (string, error) {
	dir := os.Getenv(name)
	if dir == "" {
		return "", errors.New("server getenv is not initalized")
	}
	return dir, nil
}
