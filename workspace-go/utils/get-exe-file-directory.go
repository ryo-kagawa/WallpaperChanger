package utils

import (
	"os"
	"path/filepath"
)

func GetExeFileDirectory() (string, error) {
	exeFilePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeFileDirectory := filepath.Dir(exeFilePath)
	return exeFileDirectory, nil
}
