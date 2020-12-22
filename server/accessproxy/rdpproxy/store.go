package rdpproxy

import (
	"fmt"
	"github.com/seknox/trasa/server/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/sirupsen/logrus"
)

func (s GWStore) uploadSessionLog(authlog *logs.AuthLog) error {

	tempFileDir := filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "guac")
	bucketName := "trasa-guac-logs"
	sessionID := authlog.SessionID

	loginTime := time.Unix(0, authlog.LoginTime).In(time.UTC)

	//TODO @sshahcodes

	//sudo docker exec  guacd /usr/local/guacamole/bin/guacenc -f /tmp/trasa/accessproxy/guac/%s.guac
	//here guacd is container name

	guacenc := getGuacencCmd(sessionID)

	ll, err := guacenc.CombinedOutput()
	//	logger.Debug(string(ll))
	if err != nil {
		return errors.WithMessage(err, "could not convert guac file to m4v. "+string(ll))
	} else {
		err = os.Remove(filepath.Join(tempFileDir, sessionID+".guac"))
		if err != nil {
			logrus.Errorf("could not delete guac file: %v", err)
		}
		logrus.Tracef("%s.guac file converted and deleted", sessionID)

	}

	ffmpeg := getFFMPEGcmd(tempFileDir, sessionID)
	ll, err = ffmpeg.CombinedOutput()
	//logger.Debug(string(ll))
	if err != nil {
		return errors.WithMessage(err, "could not convert m4v file to mp4. "+string(ll))
	} else {
		err = os.Remove(filepath.Join(tempFileDir, sessionID+".guac.m4v"))
		if err != nil {
			logrus.Errorf("could not delete m4v file: %v", err)
		}
		logrus.Tracef("%s.guac.m4v file converted and deleted", sessionID)

	}

	//don't use fileapth.join in object name
	objectName := fmt.Sprintf("%s/%d/%d/%d/%s.guac", authlog.OrgID, loginTime.Year(), int(loginTime.Month()), loginTime.Day(), sessionID)
	filePath := filepath.Join(tempFileDir, fmt.Sprintf("%s.mp4", sessionID))

	// Upload log file to minio
	uploadErr := logs.Store.PutIntoMinio(objectName, filePath, bucketName)
	if uploadErr != nil {
		logrus.Errorf("could not upload to minio, trying again: %v", uploadErr)
		uploadErr = logs.Store.PutIntoMinio(objectName, filePath, bucketName)
	}

	if uploadErr == nil {
		err = os.Remove(filePath)
		if err != nil {
			logrus.Errorf("could not delete mp4 file: %v", err)
		}

	}

	return uploadErr
}

func getGuacencCmd(sessionID string) *exec.Cmd {
	if os.Getenv("GUACENC_INSTALLED") == "true" {
		guacencCmdStr := fmt.Sprintf(
			"nice -n 10 /usr/local/guacamole/bin/guacenc -f /tmp/trasa/accessproxy/guac/%s.guac", sessionID)

		return exec.Command("/bin/sh", "-c", guacencCmdStr)

	}

	if runtime.GOOS == "windows" {
		guacencCmdStr := fmt.Sprintf(
			"docker.exe exec  guacd nice -n 10  /usr/local/guacamole/bin/guacenc -f /tmp/trasa/accessproxy/guac/%s.guac", sessionID)

		return exec.Command("powershell", "-c", guacencCmdStr)
	}

	guacencCmdStr := fmt.Sprintf(
		"sudo docker exec  guacd nice -n 10 /usr/local/guacamole/bin/guacenc -f /tmp/trasa/accessproxy/guac/%s.guac", sessionID)
	return exec.Command("/bin/bash", "-c", guacencCmdStr)

}

func getFFMPEGcmd(tempFileDir, sessionID string) *exec.Cmd {

	if os.Getenv("GUACENC_INSTALLED") == "true" {
		ffmpegCmdStr := fmt.Sprintf("nice -n 10 ffmpeg -i %s/%s.guac.m4v %s/%s.mp4", tempFileDir, sessionID, tempFileDir, sessionID)
		return exec.Command("/bin/bash", "-c", ffmpegCmdStr)

	}

	if runtime.GOOS == "windows" {
		ffmpegCmdStr := fmt.Sprintf(`ffmpeg.exe -i %s\%s.guac.m4v %s\%s.mp4`, tempFileDir, sessionID, tempFileDir, sessionID)
		return exec.Command("powershell", "-c", ffmpegCmdStr)

	}

	ffmpegCmdStr := fmt.Sprintf("sudo nice -n 10 ffmpeg -i %s/%s.guac.m4v %s/%s.mp4", tempFileDir, sessionID, tempFileDir, sessionID)
	return exec.Command("/bin/bash", "-c", ffmpegCmdStr)

}
