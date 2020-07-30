package my

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(con *global.State) {
	// initialize local state
	Store = MyStore{
		State: con,
	}

}

//func InitStoreMock() *userstoremocks.UserStoreMock {
//	// initialize local state
//	mockUser := new(userstoremocks.UserStoreMock)
//	Store = mockUser
//	return mockUser
//}

var Store MyAdapter

type MyStore struct {
	*global.State
}

type MyAdapter interface {
	GetUserAppsDetailsWithPolicyFromUserID(userID, orgID string) ([]models.MyServiceDetails, error)
}
