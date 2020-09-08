package utils

import (
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
		return "/var"
	}
}
