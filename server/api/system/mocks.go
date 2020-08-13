package system

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type SystemMock struct {
	mock.Mock
}

//SetGlobalSetting mock
func (s SystemMock) SetGlobalSetting(setting models.GlobalSettings) error {
	return s.Called(setting).Error(0)
}

//GetGlobalSetting mock
func (s SystemMock) GetGlobalSetting(orgID, settingType string) (models.GlobalSettings, error) {
	args := s.Called(orgID, settingType)
	return args.Get(0).(models.GlobalSettings), args.Error(1)
}

//UpdateGlobalSetting mock
func (s SystemMock) UpdateGlobalSetting(setting models.GlobalSettings) error {
	panic("implement me")
}

func (s SystemMock) getSecurityRules(orgID string) ([]models.SecurityRule, error) {
	panic("implement me")
}

func (s SystemMock) updateSecurityRule(orgID, statusstr, ruleID string) error {
	panic("implement me")
}

func (s SystemMock) storeBackupMeta(backup models.Backup) error {
	panic("implement me")
}

func (s SystemMock) getBackupMeta(backup models.Backup) (models.Backup, error) {
	panic("implement me")
}

func (s SystemMock) getBackupMetas(orgID string) ([]models.Backup, error) {
	panic("implement me")
}

//CreateSecurityRule mock
func (s SystemMock) CreateSecurityRule(rule models.SecurityRule) error {
	panic("implement me")
}

//GetSecurityRuleByName mock
func (s SystemMock) GetSecurityRuleByName(orgID, constName string) (models.SecurityRule, error) {
	panic("implement me")
}
