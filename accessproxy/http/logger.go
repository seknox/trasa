package http

import (
	"fmt"

	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func PasswordManAndLogger(r *http.Request, sessionID, csrfToken, userName string, isSSO bool, sessionRecord string) {
	if sessionRecord != "true" {
		return
	}

	var buf []byte

	addTime := fmt.Sprintf("\nTime: %s\n\n", time.Now().String())

	buf = append(buf, addTime...)

	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	directoryBuilder := fmt.Sprintf("/var/trasa/thg/logs/%s", sessionID)

	CreateDirIfNotExist(directoryBuilder)

	logPath := fmt.Sprintf("%s/%s.http-raw", directoryBuilder, sessionID)
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
}

type rawSession struct {
	Time string
	Data string
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func createFile(path string) {
	// detect if file exists
	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}

}
