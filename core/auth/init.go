package auth

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = AuthStore{state}
}

func InitStoreMock() *AuthMock {
	lmock := new(AuthMock)
	Store = lmock
	return lmock
}

var Store Adapter

type AuthStore struct {
	*global.State
}

type Adapter interface {
	GetLoginDetails(trasaID, orgDomain string) (*models.UserWithPass, error)
	Logout(sessionID string) error
}
