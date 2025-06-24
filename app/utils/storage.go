package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)

const JSON_FILE_NAME = "wcurl.json"
const DIR_NAME = ".wcurl"

type Storage struct{}

func (s *Storage) ProjectID() string {
	p, err := os.Executable()
	if err != nil {
		fmt.Println("Error:")
	}

	h := sha256.New()
	h.Write([]byte(p))
	// ph := h.Sum(nil)

	return fmt.Sprintf("%x", "foobar")
}

func (s *Storage) createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
}

func (s *Storage) GetAppPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wcurlPath := fmt.Sprintf("%s/%s", home, DIR_NAME)
	s.createDir(wcurlPath)
	return wcurlPath
}

func (s *Storage) GetAbsoluteJsonFilePath() string {
	return fmt.Sprintf("%s/%s", s.GetAppPath(), JSON_FILE_NAME)
}
