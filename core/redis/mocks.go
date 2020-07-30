package redis

import (
	"errors"
	"time"

	"github.com/seknox/trasa/consts"
	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (r RedisMock) Set(key string, expiry time.Duration, val ...string) error {

	r.TestData().Set(key, val)

	return nil
}

func (r RedisMock) Get(key string, field string) (string, error) {
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

func (r RedisMock) MGet(key string, field ...string) ([]string, error) {
	panic("implement me")
}

func (r RedisMock) Delete(key string) error {
	panic("implement me")
}

func (r RedisMock) SetVerifyIntent(key string, expiry time.Duration, intent, field, val string) error {
	panic("implement me")
}

func (r RedisMock) VerifyIntent(key string, intent consts.VerifyTokenIntent) error {
	panic("implement me")
}

func (r *RedisMock) GetSession(key string) (userID, orgID, deviceID, browserID, auth string, err error) {
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

func (r RedisMock) WaitForStatusAndGet(key, field string) (success bool, val string) {
	panic("implement me")
}

func (r RedisMock) SetHTTPGatewaySession(key, orgusr, authDataVal string, sessionRecord string) error {
	panic("implement me")
}

func (r RedisMock) GetHTTPGatewaySession(key string) (user, auth, sessionRecord string, err error) {
	panic("implement me")
}
