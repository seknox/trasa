package server_test

import (
	"github.com/seknox/trasa/tests/server/vaulttest"
	"testing"
)

func TestVault(t *testing.T) {
	vaulttest.InitVault(t)
	vaulttest.GetStatus(t)
	vaulttest.StoreVault(t)
	vaulttest.GetKey(t)
	vaulttest.StoreSecret(t)
	vaulttest.GetSecret(t)
	vaulttest.GetUpstreamCredsTest(t)
	vaulttest.DeleteSecret(t)

}
