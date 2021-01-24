package adhoc

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State) {
	Store = AdhocStore{state}
}

var Store AdhocAdapter

type AdhocStore struct {
	*global.State
}

type AdhocAdapter interface {
	Create(req models.AdhocPermission) error
	AppendSession(adhocID, sessionIDappType, orgID string) error
	GrantOrReject(req models.AdhocPermission) error
	Expire(requestID, orgID string) error
	UpdateReqSessionID(req models.AdhocPermission) error
	GetActiveReqOfUser(userID, appID, orgID string) (models.AdhocDetails, error)
	GetAll(orgID string) ([]models.AdhocDetails, error)
	GetAdhocDetail(id, orgID string) (models.AdhocDetails, error)
	GetReqsAssignedToUser(requesteeID, orgID string) ([]models.AdhocDetails, error)

	GetAdmins(orgID string) ([]models.User, error)
}
