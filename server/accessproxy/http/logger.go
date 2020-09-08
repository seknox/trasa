package http

import (
	"fmt"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"
)

func passwordManAndLogger(r *http.Request, sessionID, csrfToken, userName string, isSSO bool, sessionRecord string) error {
	if sessionRecord != "true" {
		return nil
	}

	var buf []byte

	addTime := fmt.Sprintf("\nTime: %s\n\n", time.Now().String())

	buf = append(buf, addTime...)

	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		logrus.Error(err)
		return err
	}

	directoryBuilder := fmt.Sprintf(filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "http", sessionID))

	err = createDirIfNotExist(directoryBuilder)
	if err != nil {
		return err
	}

	logPath := filepath.Join(directoryBuilder, fmt.Sprintf("%s.http-raw", sessionID))
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer file.Close()

	buf = append(buf, dump...)

	//	var modifiedBody []byte
	// if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodPut {

	// 	if isSSO == true {
	// 		modifiedBody, err := fillUnamePass(r, secrets, userName, sessionID, csrfToken)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		buf = append(buf, "\n"...)
	// 		//fmt.Println("writing modified Body; ", string(modifiedBody))
	// 		buf = append(buf, modifiedBody...)
	// 	} else {
	// 		modifiedBody, err := secretObfuscator(r, secrets, userName, sessionID, csrfToken)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		buf = append(buf, "\n"...)
	// 		//fmt.Println("writing modified Body; ", string(modifiedBody))
	// 		buf = append(buf, modifiedBody...)
	// 	}
	// }

	buf = append(buf, "\n______________________________________________________________________________________________\n"...)

	_, err = file.Write(buf)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}
