package logs

import (
	"net/http"
	"net/url"
	"time"

	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func GetVideoLog(w http.ResponseWriter, r *http.Request) {

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	ssKey := values.Get("ssKey")

	arr, err := redis.Store.MGet(ssKey, "path", "bucketName")
	if err != nil || ssKey == "" || len(arr) != 2 {
		utils.TrasaResponse(w, 200, "failed", "invalid sskey", "http video file")
		return
	}

	path := arr[0]
	bucketName := arr[1]

	// Upload the zip file
	//logger.Error(sessionID)
	// Upload log file to minio
	object, err := Store.GetFromMinio(path, bucketName)
	if err != nil {
		logrus.Error(err)
		return
	}

	http.ServeContent(w, r, "sessionID", time.Now(), object)
	return

}

func GetRawLog(w http.ResponseWriter, r *http.Request) {

	userContext := r.Context().Value("user").(models.UserContext)

	authKey := utils.GetRandomID(17)

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	sessionID := values.Get("sessionID")
	sessionType := values.Get("type")
	orgID := userContext.Org.ID
	year := values.Get("year")
	month := values.Get("month")
	day := values.Get("day")

	path, bucketName := getMinioPath(sessionID, sessionType, orgID, year, month, day)

	if sessionType != "ssh" {

		err := redis.Store.Set(authKey, time.Second*600, "path", path, "bucketName", bucketName)
		if err != nil {
			logrus.Error(err)
		}

		//fmt.Println("setting sskey ", hex.EncodeToString(authKey))
		w.Header().Set("sskey", authKey)

	}

	if sessionType == "http" {
		path, bucketName = getMinioPath(sessionID, "http-raw", orgID, year, month, day)
	}

	// Upload log file to minio
	object, err := Store.GetFromMinio(path, bucketName)
	if err != nil {
		logrus.Error(err)
		return
	}

	http.ServeContent(w, r, "sessionID", time.Now(), object)
	return
}

//getMinioPath returns minio bucket name and object name (path name)
func getMinioPath(sessionID string, sessionType, orgID, year, month, day string) (string, string) {
	bucketName := "unknown"
	objectNamePrefix := orgID + "/" + year + "/" + month + "/" + day + "/"
	objectName := objectNamePrefix + sessionID + ".session"

	//TODO add const values of session type
	if sessionType == "guac" || sessionType == "guac-ssh" || sessionType == "rdp" {
		bucketName = "trasa-guac-logs"
		objectName = objectNamePrefix + sessionID + ".guac"

		return objectName, bucketName
	}
	if sessionType == "http" {
		bucketName = "trasa-https-logs"
		objectName = objectNamePrefix + sessionID + ".mp4"

		return objectName, bucketName
	}

	if sessionType == "http-raw" {
		bucketName = "trasa-https-logs"
		objectName = objectNamePrefix + sessionID + ".http-raw"

		return objectName, bucketName
	}
	if sessionType == "ssh" {
		bucketName = "trasa-ssh-logs"
		objectName = objectNamePrefix + sessionID + ".session"

		return objectName, bucketName
	}

	if sessionType == "db" {
		bucketName = "trasa-db-logs"
		objectName = objectNamePrefix + sessionID + ".session"

		return objectName, bucketName
	}

	return objectName, bucketName

}
