package system

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = systemStore{state}
}

//InitStoreMock will init mock state of this package
func InitStoreMock() *SystemMock {
	lmock := new(SystemMock)
	Store = lmock
	return lmock
}

//Store is the package state variable which contains database connections
var Store adapter

type systemStore struct {
	*global.State
}

type adapter interface {
	// global settings
	SetGlobalSetting(setting models.GlobalSettings) error
	GetGlobalSetting(orgID, settingType string) (models.GlobalSettings, error)
	UpdateGlobalSetting(setting models.GlobalSettings) error

	// security rules
	CreateSecurityRule(rule models.SecurityRule) error
	GetSecurityRuleByName(orgID, constName string) (models.SecurityRule, error)
	getSecurityRules(orgID string) ([]models.SecurityRule, error)
	updateSecurityRule(orgID, statusstr, ruleID string) error

	// backups
	storeBackupMeta(backup models.Backup) error
	getBackupMeta(backup models.Backup) (models.Backup, error)
	getBackupMetas(orgID string) ([]models.Backup, error)
}
