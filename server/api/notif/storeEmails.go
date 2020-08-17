package notif

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/crypt"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"

	"gopkg.in/gomail.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type emailErrReport struct {
	FileName     string    `json:"fileName"`
	LineNumber   int       `json:"lineNumber"`
	Err          error     `json:"err"`
	Description  string    `json:"description"`
	FunctionName string    `json:"functionName"`
	Time         time.Time `json:"time"`
}

//Log and send error report
//This function should be called in case of severe error resulting from bug in software itself
func SendErrorReport(err error, desc string) {

	//TODO add more details to report.
	// -resources stats (memory,cpu usage)
	// -stack trace

	fileName := ""
	lineNum := -1
	functionName := ""

	fpcs := make([]uintptr, 1)
	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		logrus.Error("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		logrus.Error("MSG CALLER WAS NIL")
		//TODO possible nil pointer dereferenence
	}

	// Print the file name and line number
	fileName, lineNum = caller.FileLine(fpcs[0] - 1)

	// Print the name of the function
	functionName = caller.Name()

	tmpl := emailErrReport{
		FileName:     fileName,
		LineNumber:   lineNum,
		Err:          err,
		Description:  desc,
		FunctionName: functionName,
		Time:         time.Now(),
	}

	orgID := ""
	if global.GetConfig().Platform.Base == "private" {
		orgID = global.GetConfig().Trasa.OrgId
	} else {
		//TODO get orgID in SaaS @bhrg3se. get better explaiation first why we need to consider "SaaS" here?

	}

	logrus.Error(fmt.Sprintf(`%v  %s:%d  %s`, err, fileName, lineNum, functionName))

	sett, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_ERROR_REPORT)
	if err != nil || !sett.Status {
		return
	}

	err = Store.SendEmail(orgID, consts.EMAIL_ERR_REPORT, tmpl)
	if err != nil {
		logrus.Error(err)
	}
}

// SendEmail is single interface to send emails within TRASA operations.
// Ideally, setting values should be retrieved during start of trasacore and updated during setting update
// to remove database hit in every email sending operations. TODO @bhrg3se
func (s notifStore) SendEmail(orgID string, emailType consts.EmailType, emailTemplate interface{}) error {
	setting, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_EMAIL_CONFIG)
	if err != nil || setting.Status == false {
		return fmt.Errorf("Email setting not configured: %v", err)
	}

	// unmarshal setting value for emailIntegration struct
	var configVals models.EmailIntegrationConfig
	err = json.Unmarshal([]byte(setting.SettingValue), &configVals)
	if err != nil {
		return err
	}

	if configVals.AuthPass == "" {
		return fmt.Errorf("empty password or key not supported")
	}

	// get decrypted email password or key.
	if s.TsxvKey.State == false {
		return fmt.Errorf("encryption key is not retrieved yet.")
	}

	// get key ct from database.
	key, err := crypt.Store.GetKeyOrTokenWithKeyval(orgID, consts.GLOBAL_EMAIL_CONFIG_SECRET)
	if err != nil {
		return fmt.Errorf("eailed to retrieve cipher text.")
	}

	pt, err := utils.AESDecrypt(s.TsxvKey.Key[:], key.KeyVal)
	if err != nil {
		return fmt.Errorf("AESDecrypt. %v", err)
	}

	configVals.AuthPass = string(pt)

	// we've retrieved password or token for email, call util functino to send email.
	err = sendEmailWithConfig(configVals, emailType, emailTemplate)
	if err != nil {
		logrus.Error(err)
		//Try once more
		err = sendEmailWithConfig(configVals, emailType, emailTemplate)
		if err != nil {
			//TODO return original error err instead of this. (Maybe)

			// Because logger.Error is usually called when this function error is returned. And we need original error in logfile
			return fmt.Errorf("Failed to send email. Contact your administrator")
		}
	}
	return nil
}

// SendEmail calls specefic email sending functions depending on emailType.
func sendEmailWithConfig(emailIntegration models.EmailIntegrationConfig, emailType consts.EmailType, emailData interface{}) error {
	switch emailType {
	case consts.EMAIL_ADHOC:
		err := adhoc(emailIntegration, emailData)
		if err != nil {
			return err
		}
	case consts.EMAIL_DYNAMIC_ACCESS:
		err := dynamicAccess(emailIntegration, emailData)
		if err != nil {
			return err
		}
	case consts.EMAIL_SECURITY_ALERT:
		err := securityAlert(emailIntegration, emailData)
		if err != nil {
			return err
		}
	case consts.EMAIL_USER_CRUD:
		err := userCrud(emailIntegration, emailData)
		if err != nil {
			return err
		}

	case consts.EMAIL_ERR_REPORT:
		err := errReport(emailIntegration, emailData)
		if err != nil {
			return err
		}

	}

	return nil
}

func errReport(emailIntegration models.EmailIntegrationConfig, emailData interface{}) error {
	errReportTemplate, ok := emailData.(emailErrReport)
	if !ok {
		err := errors.New("Invalid error report template")
		return err
	}

	var t = template.New("ErrReport")
	temp, err := t.Parse(errReportTemplateBlob)

	if err != nil {
		return err
	}

	t = template.Must(temp, err)

	buf := new(bytes.Buffer)

	err = t.Execute(buf, map[string]string{
		"fileName":     errReportTemplate.FileName,
		"lineNumber":   fmt.Sprintf(`%d`, errReportTemplate.LineNumber),
		"functionName": errReportTemplate.FunctionName,
		"description":  errReportTemplate.Description,
		"err":          errReportTemplate.Err.Error(),
	})
	if err != nil {
		return err
	}

	finalHtml := buf.String()

	if emailIntegration.IntegrationType == string(consts.EMAIL_SMTP) {
		receiver := []string{"secure@seknox.com"}
		err = smtpEmail(emailIntegration, "TRASA ERROR REPORT", receiver, []string{"bhargab@seknox.com", "sakshyam@seknox.com"}, finalHtml)
		if err != nil {
			return err
		}
	} else {
		_, _, err = mailgunEmail(emailIntegration, "TRASA ERROR REPORT", "secure@seknox.com", []string{"bhargab@seknox.com", "sakshyam@seknox.com"}, finalHtml)
		if err != nil {
			return err
		}
	}
	return nil

}

// adhoc email sends adhoc access request emails.
func adhoc(emailIntegration models.EmailIntegrationConfig, emailData interface{}) error {

	var adhoc models.EmailAdhoc
	v, err := json.Marshal(emailData)
	err = json.Unmarshal(v, &adhoc)
	if err != nil {
		return err
	}

	var t = template.New("DynamicAccess")
	temp, err := t.Parse(AdhocStatus)

	if adhoc.Req {
		temp, err = t.Parse(AdhocReq)

	}

	if err != nil {
		return err
	}

	t = template.Must(temp, err)

	buf := new(bytes.Buffer)

	err = t.Execute(buf, map[string]string{"requester": adhoc.Requester, "requestee": adhoc.Requestee, "dashlink": adhoc.DashLink, "time": adhoc.Time, "app": adhoc.App, "reason": adhoc.Reason, "status": adhoc.Status})
	if err != nil {
		return err
	}

	finalHtml := buf.String()

	if emailIntegration.IntegrationType == string(consts.EMAIL_SMTP) {
		receiver := []string{adhoc.ReceiverEmail}
		err = smtpEmail(emailIntegration, adhoc.Subject, receiver, adhoc.CC, finalHtml)
		if err != nil {
			return err
		}
	} else {
		_, _, err = mailgunEmail(emailIntegration, adhoc.Subject, adhoc.ReceiverEmail, adhoc.CC, finalHtml)
		if err != nil {
			return err
		}
	}
	return nil
}

func dynamicAccess(emailIntegration models.EmailIntegrationConfig, emailData interface{}) error {

	var da models.EmailDynamicAccess
	v, err := json.Marshal(emailData)
	if err != nil {
		return err
	}
	//fmt.Println("email data: ", string(v))
	err = json.Unmarshal(v, &da)
	if err != nil {
		return err
	}

	// absPath, err := filepath.Abs("/etc/trasa/static/templates/mail/dynamicAccess.html")
	// if err != nil {
	// 	return err
	// }

	// t := template.Must(template.ParseFiles(absPath))

	var t = template.New("DynamicAccess") // *template.Template
	temp, err := t.Parse(DynamicAccess)
	if err != nil {
		logrus.Error("template error: ", err)
	}
	t = template.Must(temp, err)

	buf := new(bytes.Buffer)

	timeStr := time.Unix(da.TimeInt, 0).String()
	err = t.Execute(buf, map[string]string{"user": da.User, "hostname": da.Hostname, "appType": da.AppType, "time": timeStr})
	if err != nil {
		return err
	}

	finalHtml := buf.String()

	if emailIntegration.IntegrationType == string(consts.EMAIL_SMTP) {
		receiver := []string{da.ReceiverEmail}
		err = smtpEmail(emailIntegration, "Dyamic Service Accessed", receiver, da.CC, finalHtml)
		if err != nil {
			return err
		}
	} else {
		_, _, err = mailgunEmail(emailIntegration, "Dyamic Service Accessed", da.ReceiverEmail, da.CC, finalHtml)
		if err != nil {
			return err
		}
	}
	return nil
}

func securityAlert(emailIntegration models.EmailIntegrationConfig, emailData interface{}) error {

	var sa models.EmailSecurityAlert
	v, _ := json.Marshal(emailData)
	json.Unmarshal(v, &sa)

	var t = template.New("SecurityAlertMail") // *template.Template
	temp, err := t.Parse(SecurityAlertTemplate)
	if err != nil {
		logrus.Error("template error: ", err)
	}
	t = template.Must(temp, err)

	buf := new(bytes.Buffer)

	t.Execute(buf, map[string]string{"SecurityRuleText": sa.SecurityRuleText, "EntitynName": sa.EntityName})

	finalHtml := buf.String()

	if emailIntegration.IntegrationType == string(consts.EMAIL_SMTP) {
		receiver := []string{sa.ReceiverEmail}
		err = smtpEmail(emailIntegration, "Security Alert", receiver, sa.CC, finalHtml)
		if err != nil {
			return err
		}
	} else {
		_, _, err = mailgunEmail(emailIntegration, "Security Alert", sa.ReceiverEmail, sa.CC, finalHtml)
		if err != nil {
			return err
		}
	}
	return nil
}

func userCrud(emailIntegration models.EmailIntegrationConfig, emailData interface{}) error {
	var sa models.EmailUserCrud
	v, _ := json.Marshal(emailData)
	json.Unmarshal(v, &sa)

	subject := "TRASA Password reset"
	userInfo := sa.Username

	var t = template.New("UserCrud")
	temp, err := t.Parse(ResetPassword)

	if sa.NewM {
		temp, err = t.Parse(orgsignup)
		subject = "Welcome To Trasa"
	}

	if err != nil {
		return err
	}

	t = template.Must(temp, err)

	buf := new(bytes.Buffer)

	t.Execute(buf, map[string]string{"username": userInfo, "verifyToken": sa.VerifyUrl})

	finalHtml := buf.String()

	if emailIntegration.IntegrationType == string(consts.EMAIL_SMTP) {
		receiver := []string{sa.ReceiverEmail}
		err := smtpEmail(emailIntegration, subject, receiver, sa.CC, finalHtml)
		if err != nil {
			return err
		}
	} else {
		_, _, err := mailgunEmail(emailIntegration, subject, sa.ReceiverEmail, sa.CC, finalHtml)
		if err != nil {
			return err
		}
	}
	return nil

}

func smtpEmail(creds models.EmailIntegrationConfig, subject string, receivers []string, cc []string, emailBody string) error {
	logrus.Trace("sending smtp email")
	port, err := strconv.Atoi(creds.ServerPort)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(creds.ServerAddress, port, creds.AuthKey, creds.AuthPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: global.GetConfig().Security.InsecureSkipVerify}

	m := gomail.NewMessage()
	m.SetHeader("From", creds.AuthKey)
	m.SetHeaders(map[string][]string{
		"To": receivers,
	})
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	//d := gomail.Dialer{Host: "localhost", Port: 587}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func mailgunEmail(creds models.EmailIntegrationConfig, subject string, receiver string, cc []string, emailBody string) (string, string, error) {
	//fmt.Println("mailgun configs: ", creds.ServerAddress, creds.AuthPass, creds.SenderAddress)
	mg := mailgun.NewMailgun(creds.AuthKey, creds.AuthPass, "")
	m := mg.NewMessage(
		fmt.Sprintf("Trasa Team <%s>", creds.SenderAddress),
		//"Trasa Team <secure@trasa.io>",
		subject,
		subject,
		receiver,
	)

	m.SetHtml(emailBody)

	if len(cc) > 0 {
		for _, admin := range cc {
			m.AddCC(admin)
		}
	}

	resp, id, err := mg.Send(m)
	return resp, id, err
}

/////////////////////////////////////////
