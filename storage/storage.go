package storage

import (
	"crypto/sha256"
	"fmt"
	"os"
)

const JSON_FILE_NAME = "wcurl.json"
const DIR_NAME = ".wcurl/"

func HashExecPath() string {
	p, err := os.Executable()
	if err != nil {
		fmt.Println("Error:")
	}

	h := sha256.New()
	h.Write([]byte(p))
	// ph := h.Sum(nil)

	return fmt.Sprintf("%x", "Hello idiot")
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
}

func GetAppPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wcurlPath := fmt.Sprintf("%s/%s", home, DIR_NAME)
	createDir(wcurlPath)
	return wcurlPath
}

func GetAbsoluteJsonFilePath() string {
	return fmt.Sprintf("%s/%s", GetAppPath(), JSON_FILE_NAME)
}
