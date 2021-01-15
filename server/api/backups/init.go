package backups

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State) {
	Store = BackupStore{state}
}

var Store BackupAdapter

type BackupStore struct {
	*global.State
}

type BackupAdapter interface {
	StoreBackupMeta(backup models.Backup) error
	GetBackupMeta(backupID, orgID string) (backup models.Backup, err error)
	GetBackupMetas(orgID string) ([]models.Backup, error)
}
