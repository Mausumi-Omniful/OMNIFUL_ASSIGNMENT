package csv

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// downloadFile download file in chunks if file path is not given it automatically extracts file name from URL
func downloadFile(filepath string, url string) (filePath string, err error) {
	if filepath == "" {
		filepath = extractFileNameFromURL(url)
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file in chunks
	buf := make([]byte, 32*1024) // 32KB chunks
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return "", writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	return filepath, nil
}

func extractFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	randomNumber := generateRandomNumber()
	return fmt.Sprintf("%06d_%s", randomNumber, fileName)
}

func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000 // Generate a random number between 100000 and 999999
}
