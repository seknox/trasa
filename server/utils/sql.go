package utils

import (
	"github.com/lib/pq"
	"github.com/seknox/trasa/server/consts"
	"github.com/sirupsen/logrus"
)

//GetConstraintErrorMessage returns user readable error according to violated database constraints.
//It is used while creating and updating
func GetConstraintErrorMessage(err error) string {
	if err, ok := err.(*pq.Error); ok {
		switch err.Constraint {
		case consts.CONSTRAINT_UNIQUE_GROUPNAME:
			return "Group name already used"

		case consts.CONSTRAINT_UNIQUE_USERNAME:
			return "Username already used"
		case consts.CONSTRAINT_UNIQUE_EMAIL:
			return "Email already used"

		case consts.CONSTRAINT_UNIQUE_SERVICENAME:
			return "Service name already used"
		case consts.CONSTRAINT_UNIQUE_HOSTNAME:
			return "Hostname already used"

		default:
			logrus.Errorf("db constraint violated: %v", err)
			return "Could not create"
		}

	}

	logrus.Error(err)
	return "Could not create service"
}
