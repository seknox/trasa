package server_test

import (
	"github.com/seknox/trasa/tests/server/crud"
	"testing"
)

func TestCRUD(t *testing.T) {
	crud.CreateService(t)
	serviceID := crud.GetService(t)
	crud.UpdateService(t, serviceID)
	crud.DeleteService(t, serviceID)
}
