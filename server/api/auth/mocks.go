package auth

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type authMock struct {
	mock.Mock
}

func (a *authMock) GetLoginDetails(trasaID, orgDomain string) (*models.UserWithPass, error) {
	args := a.Called(trasaID, orgDomain)
	return args.Get(0).(*models.UserWithPass), args.Error(1)

}

func (a *authMock) Logout(sessionID string) error {
	return a.Called(sessionID).Error(0)
}
