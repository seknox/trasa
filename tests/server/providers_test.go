package server

import (
	"github.com/seknox/trasa/tests/server/idptest"
	"testing"
)

func TestProvides(t *testing.T) {
	idptest.CreateIdp(t)

	//dial tcp :636: connect: connection refused
	//TODO mock ldap server

	//idptest.UpdateIdp(t)

}
