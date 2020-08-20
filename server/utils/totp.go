package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"unicode"

	logger "github.com/sirupsen/logrus"
)

type Totp struct {
	Secret    string
	Issuer    string
	Account   string
	Algorithm string
	Digits    int
	Peroid    int
}

func GenerateTotpSecret() string {
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		logger.Error(err)
	}

	return base32.StdEncoding.EncodeToString(secret)
}

// we are adding three return string since first and last one wil be for skew values.
func CalculateTotp(dbcode string) (string, string, string) {

	finalKey, _ := decodeKey(dbcode)
	nowtime := time.Now()
	skewsub := nowtime.Add(time.Duration(-3e+10))
	skewadd := nowtime.Add(time.Duration(3e+10))

	currentCode := totp(([]byte(finalKey)), nowtime, 6)
	skewSubCode := totp(([]byte(finalKey)), skewsub, 6)
	skewAddCode := totp(([]byte(finalKey)), skewadd, 6)

	//logger.Trace(dbcode,skewSubCode, currentCode, skewAddCode)
	return fmt.Sprintf("%0*d", 6, skewSubCode), fmt.Sprintf("%0*d", 6, currentCode), fmt.Sprintf("%0*d", 6, skewAddCode)
}

func noSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}

func decodeKey(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(key))
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

func totp(key []byte, t time.Time, digits int) int {
	return hotp(key, uint64(t.UnixNano())/30e9, digits)
}
