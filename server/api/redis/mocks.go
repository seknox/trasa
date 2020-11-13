package redis

import (
	"errors"
	"time"

	"github.com/seknox/trasa/server/consts"
	"github.com/stretchr/testify/mock"
)

type redisMock struct {
	mock.Mock
}

//Set mock
func (r *redisMock) Set(key string, expiry time.Duration, val ...string) error {

	r.TestData().Set(key, val)

	return nil
}

//Get mock
func (r *redisMock) Get(key string, field string) (string, error) {
	//v:=r.s[key]

	arr, ok := r.TestData().Get(key).Data().([]string)
	if !ok {
		return "", errors.New("not found")
	}
	for i, v := range arr {
		if v == field {
			if len(arr) < i+2 {
				return "", errors.New("not found")
			}
			return arr[i+1], nil
		}
	}
	return "", errors.New("not found")
}

//MGet mock
func (r *redisMock) MGet(key string, field ...string) ([]string, error) {
	panic("implement me")
}

//Delete mock
func (r *redisMock) Delete(key string) error {
	panic("implement me")
}

//SetVerifyIntent mock
func (r *redisMock) SetVerifyIntent(key string, expiry time.Duration, intent, field, val string) error {
	panic("implement me")
}

//VerifyIntent mock
func (r redisMock) VerifyIntent(key string, intent consts.VerifyTokenIntent) error {
	panic("implement me")
}

//GetSession mock
func (r *redisMock) GetSession(key string) (userID, orgID, deviceID, browserID, auth string, err error) {
	val := r.TestData().Get(key)

	arr, ok := val.Data().([]string)
	if !ok {
		err = errors.New("key not found")
		return
	}

	for i, v := range arr {
		if v == "user" {
			if len(arr) < i+2 {
				err = errors.New("invalid array length")
				return
			}
			userID = arr[i+1]
		}
		if v == "auth" {
			if len(arr) < i+2 {
				err = errors.New("field not found")
				return
			}
			auth = arr[i+1]
		}

	}

	return
}

//WaitForStatusAndGet mock
func (r *redisMock) WaitForStatusAndGet(key, field string) (success bool, val string) {
	panic("implement me")
}

//SetHTTPAccessProxySession mock
func (r *redisMock) SetHTTPAccessProxySession(key, orgusr, authDataVal string, sessionRecord string) error {
	panic("implement me")
}

//GetHTTPAccessProxySession mock
func (r *redisMock) GetHTTPAccessProxySession(key string) (user, auth, sessionRecord string, err error) {
	panic("implement me")
}
