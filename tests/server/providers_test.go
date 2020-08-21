package server

import (
	"github.com/seknox/trasa/tests/server/idp"
	"testing"
)

func TestProvides(t *testing.T) {
	idp.CreateIdp(t)

	//dial tcp :636: connect: connection refused
	//TODO mock ldap server

	idp.UpdateIdp(t)

}
