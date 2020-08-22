package server

import (
	"github.com/seknox/trasa/tests/server/systemtest"
	"github.com/seknox/trasa/tests/server/vaulttest"
	"testing"
)

func TestSystemSettings(t *testing.T) {
	vaulttest.InitVault(t)
	systemtest.UpdateSettings(t)
	systemtest.SystemStatus(t)
	systemtest.UpdateSecurityRules(t)
	systemtest.GetSecurityRules(t)
}
