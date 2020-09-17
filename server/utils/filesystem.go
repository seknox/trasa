package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"

	logger "github.com/sirupsen/logrus"
)

//CreateDirIfNotExist creates directory if it doesn't exists
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error(err)
		}
	}
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func GetETCDir() string {
	switch runtime.GOOS {
	case "windows":
		//TODO @sshahcodes
		return ""
	default:
		return "/etc"
	}
}

func GetVarDir() string {
	switch runtime.GOOS {
	case "windows":
		return `%APPDATA%\`
	default:
		return "/var"
	}
}

func GetTmpDir() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\\Windows\TEMP`
	default:
		return "/tmp"
	}
}
