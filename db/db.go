package db

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const DEFAULT_CSV_PATH = ".uniprint/keys.csv.encrypted"
const csvPath = DEFAULT_CSV_PATH

var (
	instance *CSVManager
)

func InitializeCSV() (*CSVManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(home, csvPath)
	checkFileExistsAndCreateDirs(fullPath)
	reader := bufio.NewReader(os.Stdin)

	var pin string

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		pin, err = handleNoFileExists(fullPath, reader)
		if err != nil {
			return nil, err
		}
	} else {
		pin, err = handleFileExists(fullPath, reader)
		if err != nil {
			return nil, err
		}
	}

	csvManager, err := newCSVManager(fullPath, pin)
	if err != nil {
		return nil, err
	}
	instance = csvManager
	return csvManager, nil
}

func GetCSVManager() *CSVManager {
	if instance == nil {
		log.Panic("instance is nil. You must run InitializeCSV() before executing this function.")
	}
	return instance
}

func handleNoFileExists(fullPath string, reader *bufio.Reader) (string, error) {
	fmt.Println("No existing data found. Please set a new PIN.")

	pin, err := promptAndValidatePIN(reader, true)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"username", "password"})
	writer.Flush()

	encrypted, err := encrypt(buf.Bytes(), pin)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, encrypted, 0600); err != nil {
		return "", err
	}

	fmt.Println("New encrypted CSV file initialized.")
	return pin, nil
}

func handleFileExists(fullPath string, reader *bufio.Reader) (string, error) {
	fmt.Println("Encrypted file found. Please enter your PIN to unlock it.")

	pin, err := promptAndValidatePIN(reader, false)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	plaintext, err := decrypt(data, pin)
	if err != nil {
		return "", errors.New("failed to decrypt file. Is the PIN correct?")
	}

	fmt.Println("File successfully decrypted and loaded.")
	_ = plaintext

	return pin, nil
}
