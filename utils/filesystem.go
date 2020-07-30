package utils

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error(err)
		}
	}
}
