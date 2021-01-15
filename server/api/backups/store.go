package backups

import "github.com/seknox/trasa/server/models"

func (s BackupStore) StoreBackupMeta(backup models.Backup) error {

	_, err := s.DB.Exec(`INSERT into backups (id, org_id, name, type,created_at )
						 values($1, $2, $3, $4, $5);`, backup.BackupID, backup.OrgID, backup.BackupName, backup.BackupType, backup.CreatedAt)

	return err
}

func (s BackupStore) GetBackupMeta(backupID, orgID string) (backup models.Backup, err error) {

	err = s.DB.QueryRow(`SELECT  name,type, created_at  FROM backups WHERE id=$1 AND org_id=$2`,
		backupID, orgID).Scan(&backup.BackupName, &backup.BackupType, &backup.CreatedAt)

	return backup, err
}

func (s BackupStore) GetBackupMetas(orgID string) ([]models.Backup, error) {

	var backup models.Backup
	var backups = make([]models.Backup, 0)

	rows, err := s.DB.Query(`SELECT  id, org_id, name,type, created_at  FROM backups WHERE org_id=$1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return backups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&backup.BackupID, &backup.OrgID, &backup.BackupName, &backup.BackupType, &backup.CreatedAt)
		if err != nil {
			return backups, err
		}
		backups = append(backups, backup)
	}

	return backups, nil

}
