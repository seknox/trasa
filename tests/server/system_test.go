package server

import (
	"github.com/seknox/trasa/tests/server/system"
	"github.com/seknox/trasa/tests/server/vault"
	"testing"
)

func TestSystemSettings(t *testing.T) {
	vault.InitVault(t)
	system.UpdateSettings(t)
	system.SystemStatus(t)
}
