package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Test if the file exists
	exists := fileExists(file.Name())
	if !exists {
		t.Errorf("Expected file to exist, but it does not")
	}

	// Test if a non-existent file returns false
	exists = fileExists("nonexistentfile")
	if exists {
		t.Errorf("Expected file to not exist, but it does")
	}
}

func TestGetSha1Sum(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some data to the file
	data := []byte("test data")
	_, err = file.Write(data)
	if err != nil {
		t.Fatalf("Failed to write data to the file: %v", err)
	}

	// Calculate the expected SHA1 sum
	expectedSum := sha1.Sum(data)

	// Call the getSha1Sum function
	actualSum, err := getSha1Sum(file.Name())
	if err != nil {
		t.Fatalf("Failed to calculate SHA1 sum: %v", err)
	}

	// Compare the expected and actual SHA1 sums
	if fmt.Sprintf("%x", actualSum) != fmt.Sprintf("%x", expectedSum) {
		t.Errorf("Expected SHA1 sum %x, but got %x", expectedSum, actualSum)
	}
}
