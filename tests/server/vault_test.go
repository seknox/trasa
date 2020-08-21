package server_test

import (
	"github.com/seknox/trasa/tests/server/vault"
	"testing"
)

func TestVault(t *testing.T) {
	vault.InitVault(t)
	vault.GetStatus(t)
	vault.StoreVault(t)
	vault.GetKey(t)
	vault.StoreSecret(t)
	vault.GetSecret(t)
	vault.GetUpstreamCredsTest(t)
	vault.DeleteSecret(t)

}
