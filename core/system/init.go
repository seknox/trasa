package system

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = SystemStore{state}
}

func InitStoreMock() *SystemMock {
	lmock := new(SystemMock)
	Store = lmock
	return lmock
}

var Store SystemAdapter

type SystemStore struct {
	*global.State
}

type SystemAdapter interface {
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
