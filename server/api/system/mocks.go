package system

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type systemMock struct {
	mock.Mock
}

//SetGlobalSetting mock
func (s *systemMock) SetGlobalSetting(setting models.GlobalSettings) error {
	return s.Called(setting).Error(0)
}

//GetGlobalSetting mock
func (s *systemMock) GetGlobalSetting(orgID, settingType string) (models.GlobalSettings, error) {
	args := s.Called(orgID, settingType)
	return args.Get(0).(models.GlobalSettings), args.Error(1)
}

//UpdateGlobalSetting mock
func (s *systemMock) UpdateGlobalSetting(setting models.GlobalSettings) error {
	panic("implement me")
}

func (s *systemMock) getSecurityRules(orgID string) ([]models.SecurityRule, error) {
	panic("implement me")
}

func (s *systemMock) updateSecurityRule(orgID, statusstr, ruleID string) error {
	panic("implement me")
}

func (s *systemMock) storeBackupMeta(backup models.Backup) error {
	panic("implement me")
}

func (s *systemMock) getBackupMeta(backupID, orgID string) (backup models.Backup, err error) {
	panic("implement me")
}

func (s *systemMock) getBackupMetas(orgID string) ([]models.Backup, error) {
	panic("implement me")
}

//CreateSecurityRule mock
func (s *systemMock) CreateSecurityRule(rule models.SecurityRule) error {
	panic("implement me")
}

//GetSecurityRuleByName mock
func (s *systemMock) GetSecurityRuleByName(orgID, constName string) (models.SecurityRule, error) {
	panic("implement me")
}
