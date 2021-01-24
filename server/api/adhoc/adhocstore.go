package adhoc

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/seknox/trasa/server/models"
)

// StoreAdhocReq stores adhoq application access request. pq.Array(req.SessionID)
func (s AdhocStore) Create(req models.AdhocPermission) error {

	_, err := s.DB.Exec(`INSERT into adhoc_perms (id, requester_id, service_id, org_id, requestee_id, request_text, requested_on, is_authorized, authorized_on, authorized_period, authorized_policy, session_id,is_expired)
						 values($1, $2, $3, $4, $5, $6, $7, $8 ,$9, $10,$11, $12,$13);`, req.RequestID, req.RequesterID, req.ServiceID, req.OrgID, req.RequesteeID, req.RequestTxt, req.RequestedOn, req.IsAuthorized, req.AuthorizedOn, req.AuthorizedPeriod, req.AuthorizedPolicy, pq.Array(req.SessionID), req.IsExpired)

	return err
}

func (s AdhocStore) AppendSession(adhocID, sessionIDappType, orgID string) error {
	_, err := s.DB.Exec(`UPDATE adhoc_perms SET session_id=array_append((SELECT session_id from adhoc_perms WHERE id=$2 AND org_id=$3),$1) WHERE id=$2 AND org_id=$3`,
		sessionIDappType, adhocID, orgID)

	return err
}

// StoreAdhocReq stores adhoq application access request. pq.Array(req.SessionID)
func (s AdhocStore) GrantOrReject(req models.AdhocPermission) error {

	_, err := s.DB.Exec(`UPDATE adhoc_perms set is_authorized=$3, authorized_on=$4, authorized_period=$5, authorized_policy=$6,is_expired=$7 
									WHERE id=$1 AND org_id=$2;`,
		req.RequestID, req.OrgID, req.IsAuthorized, req.AuthorizedOn, req.AuthorizedPeriod, req.AuthorizedPolicy, req.IsExpired)

	return err
}

// StoreAdhocReq stores adhoq application access request. pq.Array(req.SessionID)
func (s AdhocStore) Expire(requestID, orgID string) error {

	_, err := s.DB.Exec(`UPDATE adhoc_perms set is_expired= $3 WHERE id=$1 AND org_id=$2;`, requestID, orgID, true)

	return err
}

//TODO is this used??
func (s AdhocStore) UpdateReqSessionID(req models.AdhocPermission) error {

	_, err := s.DB.Exec(`UPDATE adhoc_perms set is_expired=$3, session_id=$4 WHERE request_id=$1 AND org_id=$2;`,
		req.RequestID, req.OrgID, req.IsExpired, req.SessionID)

	return err
}

func (s AdhocStore) GetActiveReqOfUser(userID, appID, orgID string) (models.AdhocDetails, error) {
	var req models.AdhocDetails

	err := s.DB.QueryRow(`select adhoc_perms.id,requester.email,requestee.email, requester_id, service_id,services.name, requestee_id, request_text, requested_on, is_authorized, authorized_on, authorized_period, authorized_policy, session_id
									from adhoc_perms
									JOIN users requester ON adhoc_perms.requester_id=requester.id
									JOIN users requestee ON adhoc_perms.requestee_id=requestee.id
									JOIN services ON adhoc_perms.service_id=services.id
								where requester_id=$1 AND service_id=$2 AND adhoc_perms.org_id=$3 AND is_expired=$4;`,
		userID, appID, orgID, false).
		Scan(&req.RequestID, &req.RequesterEmail, &req.RequesteeEmail, &req.RequesterID, &req.ServiceID, &req.ServiceName, &req.RequesteeID, &req.RequestTxt, &req.RequestedOn, &req.IsAuthorized, &req.AuthorizedOn, &req.AuthorizedPeriod, &req.AuthorizedPolicy, pq.Array(&req.SessionID))

	return req, err
}

func (s AdhocStore) GetAdhocDetail(id, orgID string) (models.AdhocDetails, error) {
	var req models.AdhocDetails

	err := s.DB.QueryRow(`select adhoc_perms.id,requester.email,requestee.email, requester_id, service_id,services.name,services.type, requestee_id, request_text, requested_on, is_authorized, authorized_on, authorized_period, authorized_policy, session_id 
									from adhoc_perms
										JOIN users requester ON adhoc_perms.requester_id=requester.id
										JOIN users requestee ON adhoc_perms.requestee_id=requestee.id
										JOIN services ON adhoc_perms.service_id=services.id

									WHERE  adhoc_perms.id=$1 AND adhoc_perms.org_id =$2;`, id, orgID).Scan(&req.RequestID, &req.RequesterEmail, &req.RequesteeEmail, &req.RequesterID, &req.ServiceID, &req.ServiceName, &req.ServiceType, &req.RequesteeID, &req.RequestTxt, &req.RequestedOn, &req.IsAuthorized, &req.AuthorizedOn, &req.AuthorizedPeriod, &req.AuthorizedPolicy, pq.Array(&req.SessionID))
	return req, err
}

// GetAll should only return expired row.
func (s AdhocStore) GetAll(orgID string) ([]models.AdhocDetails, error) {

	var reqs []models.AdhocDetails

	rows, err := s.DB.Query(`select adhoc_perms.id,requester.email,requestee.email, requester_id, service_id,services.name,services.type, requestee_id, request_text, requested_on, is_authorized, authorized_on, authorized_period, authorized_policy, session_id 
									from adhoc_perms
										JOIN users requester ON adhoc_perms.requester_id=requester.id
										JOIN users requestee ON adhoc_perms.requestee_id=requestee.id
										JOIN services ON adhoc_perms.service_id=services.id

									WHERE adhoc_perms.is_expired =$1 AND adhoc_perms.org_id =$2 ORDER BY requested_on DESC;`, true, orgID)

	if err != nil {
		return reqs, err

	}
	defer rows.Close()

	for rows.Next() {
		var req models.AdhocDetails
		err = rows.Scan(&req.RequestID, &req.RequesterEmail, &req.RequesteeEmail, &req.RequesterID, &req.ServiceID, &req.ServiceName, &req.ServiceType, &req.RequesteeID, &req.RequestTxt, &req.RequestedOn, &req.IsAuthorized, &req.AuthorizedOn, &req.AuthorizedPeriod, &req.AuthorizedPolicy, pq.Array(&req.SessionID))
		if err != nil {
			return reqs, err

		}
		reqs = append(reqs, req)
	}

	return reqs, nil
}

//  select user_id, service_id from inapp_notifs join adhoc_perms on inapp_notifs.entity_id=adhoc_perms.request_id where user_id='sakshyam@seknox.com' AND is_resolved=false;
// StoreAdhocReq stores adhoq application access request.
func (s AdhocStore) GetReqsAssignedToUser(requesteeID, orgID string) ([]models.AdhocDetails, error) {
	var req models.AdhocDetails
	var reqs []models.AdhocDetails

	rows, err := s.DB.Query(`select adhoc_perms.id,requester.email,requestee.email, requester_id, service_id,services.name, requestee_id, request_text, requested_on, is_authorized, authorized_on, authorized_period, authorized_policy, session_id 
									from inapp_notifs 
									    join adhoc_perms on inapp_notifs.emitter_id=adhoc_perms.id
									    JOIN users requester ON adhoc_perms.requester_id=requester.id
										JOIN users requestee ON adhoc_perms.requestee_id=requestee.id
									JOIN services ON adhoc_perms.service_id=services.id

									where user_id=$1 AND is_resolved=$2 AND inapp_notifs.org_id=$3 ORDER BY requested_on DESC;`, requesteeID, false, orgID)
	if err != nil {
		return reqs, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&req.RequestID, &req.RequesterEmail, &req.RequesteeEmail, &req.RequesterID, &req.ServiceID, &req.ServiceName, &req.RequesteeID, &req.RequestTxt, &req.RequestedOn, &req.IsAuthorized, &req.AuthorizedOn, &req.AuthorizedPeriod, &req.AuthorizedPolicy, pq.Array(&req.SessionID))
		if err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (s AdhocStore) GetAdmins(orgID string) ([]models.User, error) {

	var admins []models.User

	rows, err := s.DB.Query(`SELECT  id,first_name, 
	last_name, email FROM users WHERE user_role = $1 and org_id = $2`, "orgAdmin", orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return admins, err
		}
		admins = append(admins, user)
	}

	return admins, nil

}

func (s AdhocStore) AppendAdhocSession(adhocID, sessionIDappType, orgID string) error {
	_, err := s.DB.Exec(`UPDATE adhoc_perms SET session_id=array_append((SELECT session_id from adhoc_perms WHERE request_id=$2 AND org_id=$3),$1) WHERE id=$2 AND org_id=$3`, sessionIDappType, adhocID, orgID)

	if err != nil {
		return fmt.Errorf("%d", 1)

	}
	return nil
}
