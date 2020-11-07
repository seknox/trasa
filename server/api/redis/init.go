package redis

import (
	"time"

	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = redisStore{State: state}
}

//InitStoreMock will init mock state of this package
func InitStoreMock() *redisMock {
	m := new(redisMock)
	Store = m
	return m
}

//Store is the package state variable which contains database connections
var Store adapter

type redisStore struct {
	*global.State
}

type adapter interface {
	Set(key string, expiry time.Duration, val ...string) error
	Get(key string, field string) (string, error)
	MGet(key string, field ...string) ([]string, error)
	Delete(key string) error

	SetVerifyIntent(key string, expiry time.Duration, intent, field, val string) error
	VerifyIntent(key string, intent consts.VerifyTokenIntent) error

	GetSession(key string) (userID, orgID, deviceID, browserID, auth string, err error)

	WaitForStatusAndGet(key, field string) (success bool, val string)
	SetHTTPAccessProxySession(key, orgusr, authDataVal string, sessionRecord string) error
	GetHTTPAccessProxySession(key string) (user, auth, sessionRecord string, err error)
}
