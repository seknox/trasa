package utils

import (
	"github.com/pkg/errors"
	"strings"
)

//NormalizeString trims spaces and convert into lowercase
func NormalizeString(s string) string {
	s = strings.TrimSpace(s)

	s = strings.ToLower(s)
	return s
}

// DomainFromEmail deduce domain name from email address. this domain is used to create subdomain on trasa-cloud
func DomainFromEmail(email string) string {

	//// get org name
	//index := strings.Index(email, "@")
	//lastindex := strings.LastIndex(email, ".")
	////return "", email[(index + 1):lastindex]

	// get domain name
	pos := strings.LastIndex(email, "@")
	if pos == -1 {
		return ""
	}
	val := pos + len("@")
	if val >= len(email) {
		return ""
	}
	return email[val:len(email)]

}

//ArrayContainsString check if an string array contains a string
func ArrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//ToStringArr casts interface array into string array
func ToStringArr(vals []interface{}) ([]string, error) {
	var resultStrArr = make([]string, 0)
	for _, v := range vals {
		str, ok := v.(string)
		if !ok {
			return resultStrArr, errors.Errorf("could not cast value to string")
		}
		resultStrArr = append(resultStrArr, str)
	}

	return resultStrArr, nil
}
