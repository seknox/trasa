package server_test

import (
	"github.com/seknox/trasa/tests/server/crudtest"
	"testing"
)

func TestCRUD(t *testing.T) {
	crudtest.CreateService(t)
	serviceID := crudtest.GetService(t)
	crudtest.UpdateService(t, serviceID)
	crudtest.DeleteService(t, serviceID)
}
