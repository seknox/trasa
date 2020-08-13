package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
)

//GetRandomString returns random string
func GetRandomString(length int) string {
	val := make([]byte, length)
	_, _ = rand.Read(val)
	// if err != nil {

	// }

	return hex.EncodeToString(val)

}

// GetRandomBytes returns crypto rand bytes
func GetRandomBytes(length int) []byte {
	val := make([]byte, length)
	_, _ = rand.Read(val)

	return val

}

//if uuid.Newv4() gets error, it panics
func GetUUID() string {
	uuID, err := uuid.NewV4()
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return uuID.String()
}
