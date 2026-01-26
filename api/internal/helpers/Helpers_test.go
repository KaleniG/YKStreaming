package helpers

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomToken(t *testing.T) {
	const tokenLength uint32 = 10
	token, err := GenerateRandomToken(tokenLength)

	assert.Nil(t, err)
	assert.Equal(t, tokenLength, uint32(len([]byte(token))))
}

func TestFileExists(t *testing.T) {
	// Setup
	randomToken, err := GenerateRandomToken(10)
	assert.Nil(t, err)
	tmpDir := t.TempDir()
	filename := tmpDir + "/" + randomToken + ".txt"

	// Test 1
	{
		assert.False(t, FileExists(filename))
	}

	// Test 2
	{
		file, err := os.Create(filename)
		assert.Nil(t, err)
		defer file.Close()

		assert.True(t, FileExists(filename))
	}
}

func TestFindFilesContaining(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	tmpDir2 := t.TempDir()

	files := []string{
		"apple.txt",
		"banana.txt",
		"grape.log",
		"pineapple.txt",
	}
	for _, name := range files {
		f, err := os.Create(tmpDir + "/" + name)
		assert.Nil(t, err)
		f.Close()
	}

	// Test 1: Find files containing "apple"
	{
		matched, err := FindFilesContaining(tmpDir, "apple")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"apple.txt", "pineapple.txt"}, matched)
	}

	// Test 2: Find files containing "banana"
	{
		matched, err := FindFilesContaining(tmpDir, "banana")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"banana.txt"}, matched)
	}

	// Test 3: Find files containing "orange" (none)
	{
		matched, err := FindFilesContaining(tmpDir, "orange")
		assert.Nil(t, err)
		assert.Empty(t, matched)
	}

	// Test 4: Find files in wrong dir
	{
		matched, err := FindFilesContaining(tmpDir2, "apple")
		assert.Nil(t, err)
		assert.Empty(t, matched)
	}
}

func TestGetEnvDir(t *testing.T) {
	// Test 1: Before init call
	{
		_, err := GetEnvDir("THUMBNAILS_DIR")
		assert.Error(t, err)
	}

	// Setup
	err := godotenv.Load("./../../.env")
	assert.Nil(t, err)

	// Test 2: After init call
	{
		_, err := GetEnvDir("THUMBNAILS_DIR")
		assert.Nil(t, err)
	}

	// Test 3: After init call
	{
		randomToken, err := GenerateRandomToken(10)
		assert.Nil(t, err)

		_, err = GetEnvDir(randomToken)
		assert.Error(t, err)
	}
}
