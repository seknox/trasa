package redis

import (
	"time"

	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/global"
)

func InitStore(state *global.State) {
	Store = RedisStore{State: state}
}

func InitStoreMock() *RedisMock {
	m := new(RedisMock)
	Store = m
	return m
}

var Store RedisAdapter

type RedisStore struct {
	*global.State
}

type RedisAdapter interface {
	Set(key string, expiry time.Duration, val ...string) error
	Get(key string, field string) (string, error)
	MGet(key string, field ...string) ([]string, error)
	Delete(key string) error

	SetVerifyIntent(key string, expiry time.Duration, intent, field, val string) error
	VerifyIntent(key string, intent consts.VerifyTokenIntent) error

	GetSession(key string) (userID, orgID, deviceID, browserID, auth string, err error)

	WaitForStatusAndGet(key, field string) (success bool, val string)
	SetHTTPGatewaySession(key, orgusr, authDataVal string, sessionRecord string) error
	GetHTTPGatewaySession(key string) (user, auth, sessionRecord string, err error)
}
