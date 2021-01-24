package system

import (
	"bufio"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mholt/archiver"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func TakeBackupNow(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var backup models.Backup
	// get unique backup id
	backup.BackupID = utils.GetUUID()
	backup.OrgID = userContext.User.OrgID

	orgDetails, err := orgs.Store.Get(userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to org data", "Backup not taken")
		return
	}

	loc, err := time.LoadLocation(orgDetails.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Invalid timezone", "Backup not taken")
		return
	}

	// get unique backup name
	backup.BackupName = fmt.Sprintf("trasa-backup-%s", time.Now().In(loc).Format(time.RFC3339))

	trasaBackupDir := fmt.Sprintf("/var/trasa/backup/%s", backup.BackupName)

	back, err := takeSysBackup(trasaBackupDir, backup)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to take backup", "Backup not taken")
		return
	}

	archive(trasaBackupDir, fmt.Sprintf("%s.zip", trasaBackupDir))

	// store backup metadata in database
	err = Store.storeBackupMeta(back)
	if err != nil {
		logrus.Errorf("failed to store backup meta: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to store backup meta", "Backup not taken")
		return
	}

	resp, err := Store.getBackupMetas(userContext.User.OrgID)
	if err != nil {
		logrus.Errorf("error retrieving backup metas: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch backup meta", "Backup not taken")
		return
	}

	utils.TrasaResponse(w, 200, "success", "backup created", "Backup taken", resp)
}

// takeSysBackup takes current snapshot of cockroachdb.
func takeSysBackup(trasaBackupDir string, backup models.Backup) (models.Backup, error) {

	backup.BackupType = "SYSTEM"

	backup.CreatedAt = time.Now().Unix()

	err := backupCRDB(trasaBackupDir)
	if err != nil {
		logrus.Errorf("error in takeSysBackup: %v", err)
		return backup, err
	}

	return backup, nil

}

func backupCRDB(trasaBackupDir string) error {

	// get backup directory

	err := os.MkdirAll(trasaBackupDir, 0655)
	if err != nil {
		return err
	}

	// create

	var cmd *exec.Cmd
	if global.GetConfig().Database.Dbtype == "postgres" {
		cmd = exec.Command(
			"/bin/sh",
			"-c",
			fmt.Sprintf("env PGPASSWORD=%s  pg_dump --inserts --host=%s --port=%s --username=%s --dbname=trasadb", global.GetConfig().Database.Dbpass, global.GetConfig().Database.Server, global.GetConfig().Database.Port, global.GetConfig().Database.Dbuser))
	} else if global.GetConfig().Database.Dbtype == "cockroachdb" {
		cmd = exec.Command("cockroach", "dump", "trasadb", "--insecure")
		if global.GetConfig().Database.Sslenabled {
			logrus.Trace(global.GetConfig().Database.Sslenabled)
			cmd = exec.Command("cockroach", "dump", "trasadb", "--certs-dir=/etc/trasa/certs",
				fmt.Sprintf("--host=%s:%s", global.GetConfig().Database.Server, global.GetConfig().Database.Port))
		}

	} else {
		return fmt.Errorf("invalid database type specified in config: %s", global.GetConfig().Database.Dbtype)
	}

	cmd.Dir = trasaBackupDir

	// open the out file for writing
	outfile, err := os.Create(fmt.Sprintf("%s/%s", trasaBackupDir, "cockroach-back.sql"))
	if err != nil {
		return err
	}
	defer outfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = cmd.Start()
	if err != nil {
		return err
	}

	go io.Copy(writer, stdoutPipe)
	return cmd.Wait()

}

func GetBackups(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	resp, err := Store.getBackupMetas(userContext.User.OrgID)
	if err != nil {
		logrus.Errorf("error retrieving backup metas %v: ", err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch backup meta", "GetBackups")
		return
	}

	latestOnly := r.URL.Query().Get("latest")
	if latestOnly == "true" {
		var latestBackup models.Backup
		if len(resp) > 0 {
			latestBackup = resp[0]
		}
		utils.TrasaResponse(w, 200, "success", "Backups fetched", "GetBackups", latestBackup, len(resp))
		return
	}

	utils.TrasaResponse(w, 200, "success", "Backups fetched", "GetBackups", resp)

}

func archive(archivePath, outputName string) {

	err := archiver.Archive([]string{archivePath}, outputName)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func DownloadBackupFile(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	backupID := chi.URLParam(r, "backupid")

	// get backup detail from backup id
	backup, err := Store.getBackupMeta(backupID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch backup meta", "GetBackups")
		return
	}

	// get backup directory
	trasaBackupFile := fmt.Sprintf("/var/trasa/backup/%s.zip", backup.BackupName)
	logrus.Tracef(`Backup file is %s`, trasaBackupFile)

	attachmentName := fmt.Sprintf("attachment; filename=%s.zip", backup.BackupName)

	w.Header().Set("Content-Disposition", attachmentName)
	w.Header().Set("Content-Type", "application/zip")

	http.ServeFile(w, r, trasaBackupFile)
}
