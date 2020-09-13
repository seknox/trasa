package rdpproxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

func (s GWStore) CheckPolicy(params *models.ConnectionParams, policy *models.Policy, adhoc bool) (bool, consts.FailedReason) {
	return s.checkPolicyFunc(params, policy, adhoc)
}

func (s GWStore) uploadSessionLog(authlog *logs.AuthLog) error {

	tempFileDir := "/tmp/trasa/accessproxy/guac"
	bucketName := "trasa-guac-logs"
	sessionID := authlog.SessionID
	logrus.Debugf("sessionID is %s", sessionID)

	loginTime := time.Unix(0, authlog.LoginTime)

	guacencCmdStr := fmt.Sprintf("sudo docker exec  guacd /usr/local/guacamole/bin/guacenc -f /tmp/trasa/accessproxy/guac/%s.guac", sessionID)
	guacenc := exec.Command("/bin/bash", "-c", guacencCmdStr)
	ll, err := guacenc.CombinedOutput()
	//	logger.Debug(string(ll))
	if err != nil {
		return errors.WithMessage(err, "could not convert guac file to m4v. "+string(ll))
	} else {
		err = os.Remove(filepath.Join(tempFileDir, authlog.SessionID+".guac"))
		if err != nil {
			logrus.Errorf("could not delete mp4 file: %v", err)
		}
		logrus.Tracef("%s.guac file converted and deleted", sessionID)

	}

	ffmpegCmdStr := fmt.Sprintf("sudo ffmpeg -i %s/%s.guac.m4v %s/%s.mp4", tempFileDir, sessionID, tempFileDir, sessionID)
	ffmpeg := exec.Command("/bin/bash", "-c", ffmpegCmdStr)
	ll, err = ffmpeg.CombinedOutput()
	//logger.Debug(string(ll))
	if err != nil {
		return errors.WithMessage(err, "could not convert m4v file to mp4. "+string(ll))
	} else {
		err = os.Remove(filepath.Join(tempFileDir, authlog.SessionID+".guac.m4v"))
		if err != nil {
			logrus.Errorf("could not delete mp4 file: %v", err)
		}
		logrus.Tracef("%s.guacamole.mp4 file converted and deleted", sessionID)

	}

	objectName := fmt.Sprintf("%s/%d/%d/%d/%s.guac", authlog.OrgID, loginTime.Year(), int(loginTime.Month()), loginTime.Day(), sessionID)
	filePath := fmt.Sprintf("%s/%s.mp4", tempFileDir, sessionID)

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
