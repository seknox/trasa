package utils

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
)

func MarshallStructStr(s interface{}) string {
	d, err := json.Marshal(s)
	if err != nil {
		logger.Debug(err)
	}
	return string(d)
}

//MarshallStructByte marshalls interface into bytes ignoring errors
func MarshallStructByte(s interface{}) []byte {
	d, err := json.Marshal(s)
	if err != nil {
		logger.Debug(err)
	}
	return d
}
