package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strings"
)

func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

// TO EXECUTE ONLY AFTER SERVER INIT
func GetEnvDir(name string) string {
	dir := os.Getenv(name)
	if dir == "" {
		log.Fatal(name + " is missing")
	}
	return dir
}
