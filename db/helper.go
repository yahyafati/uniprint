package db

import (
	"os"
	"path/filepath"
)

func checkFileExistsAndCreateDirs(filePath string) (bool, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return false, err
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
