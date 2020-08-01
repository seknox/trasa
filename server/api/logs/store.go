package logs

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	"github.com/minio/minio-go"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
)

//TODO remove unnecessary fields
const logparams = `id ,
    session_id ,
    access_device_id ,
    tfa_device_id ,
    service_id  ,
    service_name ,
	service_type,
    email ,
    failed_reason ,
    geo_location ,
    login_time ,
    logout_time ,
    org_id   ,
    server_ip ,
    session_duration ,
    status ,
    user_agent ,
    user_id ,
    user_ip ,
    privilege,
	guests,
	recorded_session`

const inappptrailparams = `
	id ,
	client_ip ,
	user_agent ,
	email ,
	event_time ,
	org_id ,
	status ,
	user_id,
	description
`

func (s LogStore) LogSignup(signup *models.InitSignup) error {

	var logData models.SignupLog

	logData.Company = signup.Company
	logData.Country = signup.Country
	logData.Email = signup.Email
	logData.FirstName = signup.FirstName
	logData.LastName = signup.LastName
	logData.PhoneNumber = signup.PhoneNumber
	logData.Reference = signup.Reference
	//logData.SignupTime = time.Now().Format(time.RFC3339)
	logData.SignupTime = time.Now().Unix()

	_, err := s.DB.Exec(`INSERT INTO signup_logs (company, country, email, first_name, last_name, phone_number, reference, signup_time) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)

	return err

}

func (s LogStore) LogLogin(log *AuthLog, reason consts.FailedReason, status bool) error {

	log.Status = status
	log.LogoutTime = time.Now().UnixNano()
	log.FailedReason = reason

	if (!log.Status) && (log.FailedReason == "") {
		log.FailedReason = consts.REASON_UNKNOWN
	}

	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO auth_logs (%s) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`, logparams),
		log.EventID,
		log.SessionID,
		log.AccessDeviceID,
		log.TfaDeviceID,
		log.ServiceID,
		log.ServiceName,
		log.ServiceType,
		log.Email,
		log.FailedReason,
		log.GeoLocation,
		log.LoginTime,
		log.LogoutTime,
		log.OrgID,
		log.ServerIP,
		log.LogoutTime-log.LoginTime,
		log.Status,
		log.UserAgent,
		log.UserID,
		log.UserIP,
		log.Privilege,
		pq.Array(log.Guests),
		log.SessionRecord,
	)
	return err
}

func (s LogStore) GetLoginEvents(entityType, entityID, orgID string) (logEvents []AuthLog, err error) {

	sqlStr := ""

	if entityType == "user" {
		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND user_id=$2 ORDER BY login_time DESC LIMIT 100`, logparams)

	} else if entityType == "service" {
		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND service_id=$2 ORDER BY login_time DESC LIMIT 100`, logparams)
	} else if entityType == "org" {
		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND org_id=$2 ORDER BY login_time DESC LIMIT 100`, logparams)
	} else {
		return logEvents, nil
	}

	return querySQLAuth(s.DB, sqlStr, orgID, entityID)

}

func (s LogStore) GetLoginEventsByPage(entityType, entityID, orgID string, page int, size int, dateFrom, dateTo int64) ([]AuthLog, error) {

	sqlStr := ``

	if entityType == "user" {
		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND user_id=$2 AND login_time>$3 AND login_time<$4 ORDER BY login_time DESC LIMIT $5 OFFSET $6`, logparams)
		return querySQLAuth(s.DB, sqlStr, orgID, entityID, dateFrom, dateTo, size, page)

	} else if entityType == "service" {

		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND service_id=$2 AND login_time>$3 AND login_time<$4 ORDER BY login_time DESC LIMIT $5 OFFSET $6`, logparams)
		return querySQLAuth(s.DB, sqlStr, orgID, entityID, dateFrom, dateTo, size, page)

	} else {
		sqlStr = fmt.Sprintf(`SELECT %s FROM auth_logs WHERE org_id=$1 AND login_time>$2 AND login_time<$3 ORDER BY login_time DESC LIMIT $4 OFFSET $5`, logparams)
		return querySQLAuth(s.DB, sqlStr, orgID, dateFrom, dateTo, size, page)

	}

}

func querySQLAuth(conn *sql.DB, sqlStr string, arg ...interface{}) ([]AuthLog, error) {
	var logEvents []AuthLog = make([]AuthLog, 0)
	rows, err := conn.Query(sqlStr, arg...)
	if err != nil {
		return logEvents, err
	}
	// Iterate through results
	for rows.Next() {
		var log AuthLog
		err := rows.Scan(&log.EventID,
			&log.SessionID,
			&log.AccessDeviceID,
			&log.TfaDeviceID,
			&log.ServiceID,
			&log.ServiceName,
			&log.ServiceType,
			&log.Email,
			&log.FailedReason,
			&log.GeoLocation,
			&log.LoginTime,
			&log.LogoutTime,
			&log.OrgID,
			&log.ServerIP,
			&log.SessionDuration,
			&log.Status,
			&log.UserAgent,
			&log.UserID,
			&log.UserIP,
			&log.Privilege,
			pq.Array(&log.Guests),
			&log.SessionRecord,
		)
		if err != nil {
			return logEvents, err
		}

		logEvents = append(logEvents, log)
	}

	//mar, _ := json.Marshal(logEvents)
	return logEvents, nil
}

func (s LogStore) GetOrgInAppTrails(orgID string, page int, size int, dateFrom, dateTo int64) ([]models.InAppTrail, error) {

	var logEvents = make([]models.InAppTrail, 0)

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select(strings.Split(inappptrailparams, ",")...)
	sb.From("inapp_trails")
	sb.Where(sb.Equal("org_id", orgID))
	sb.OrderBy("event_time").Desc()
	sb.Limit(size)
	sb.Offset(page)

	if dateTo > 0 {
		sb.LessEqualThan("event_time", dateTo)
	}
	if dateFrom > 0 {
		sb.GreaterEqualThan("event_time", dateFrom)
	}

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)

	logEvents, err := querySQLInappTrail(s.DB, sqlStr, args...)
	if err != nil {
		return logEvents, err
	}

	return logEvents, nil
}

func querySQLInappTrail(conn *sql.DB, sqlStr string, arg ...interface{}) ([]models.InAppTrail, error) {
	var logEvents []models.InAppTrail = make([]models.InAppTrail, 0)
	rows, err := conn.Query(sqlStr, arg...)
	if err != nil {
		return logEvents, err
	}
	//	var ignore string
	// Iterate through results
	for rows.Next() {
		var log models.InAppTrail

		err := rows.Scan(
			&log.EventID,
			&log.ClientIP,
			&log.UserAgent,
			&log.Email,
			&log.EventTime,
			&log.OrgID,
			&log.Status,
			&log.UserID,
			&log.Description,
		)
		if err != nil {
			logrus.Error(err)
			return logEvents, err
		}

		logEvents = append(logEvents, log)
	}

	//mar, _ := json.Marshal(logEvents)
	return logEvents, nil
}

func (s LogStore) LogInAppTrail(ip, userAgent, description string, user *models.User, status bool) error {

	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO inapp_trails (%s) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, inappptrailparams),
		utils.GetRandomID(5),
		ip,
		userAgent,
		user.Email,
		time.Now().UnixNano(),
		user.OrgID,
		status,
		user.ID,
		description,
	)
	if err != nil {
		logrus.Error(err)
	}

	return err
}

func (s LogStore) GetFromMinio(path, bucketName string) (*minio.Object, error) {
	// Download log file to minio
	object, err := s.MinioClient.GetObject(bucketName, path, minio.GetObjectOptions{})

	return object, err

}

func (s LogStore) PutIntoMinio(path, filepath, bucketName string) error {
	// Download log file to minio
	_, err := s.MinioClient.FPutObject(bucketName, path, filepath, minio.PutObjectOptions{})
	return err

}

// UploadHTTPLogToMinio uploads http txt and video log to minio
func (s LogStore) UploadHTTPLogToMinio(file *os.File, login AuthLog) error {

	bucketName := "trasa-https-logs"
	filePath := file.Name()
	loginTime := time.Unix(0, login.LoginTime)
	objectNamePrefix := login.OrgID + "/" + strconv.Itoa(loginTime.Year()) + "/" + strconv.Itoa(int(loginTime.Month())) + "/" + strconv.Itoa(loginTime.Day()) + "/"

	objectName := objectNamePrefix + filepath.Base(file.Name())

	contentType := "text/plain"

	// Upload log file to minio
	n, err := s.MinioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logrus.Debug(err)
		return err
	}

	logrus.Tracef("Successfully uploaded %s of size %d to minio \n", objectName, n)
	return nil
}
