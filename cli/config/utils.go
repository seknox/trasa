package config

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Context struct {
	TRASA_URL       string
	TRASA_HOST      string
	OS_TYPE         string
	SSH_CLIENT_NAME string
	TRASA_ORG_ID    string
	KEYS_DIR_PATH   string
	TRASA_ID        string
	NEW_TRASA_ID    *bool
	SSH_USERNAME    *string
	RAW_ARGS        *string
	OS_USERNAME     string
	U_ID            int
	G_ID            int
}

func Pipe(c1, c2 *exec.Cmd) {
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
}

func NormalizeString(s string) string {
	s = strings.TrimSpace(s)

	s = strings.ToLower(s)
	return s
}

//func SubmitAndDownload(data interface{}) (*models.TrasaResponse,error) {
//
//	url:=TRASA_HOSTNAME+"/idp/login/deviceAgent"
//	//url:="http://localhost:3339/"+path
//
//
//	mars,err:=json.Marshal(data)
//	if(err!=nil){
//		return nil,err
//	}
//
//
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(mars))
//	if err != nil {
//		return nil,err
//	}
//

//	//req.Header.Set("X-CSRF", csrfToken)
//
//
//	fmt.Printf("request sent was: %s\n", req.RequestURI)
//
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{Transport: tr}
//	resp, err := client.Do(req)
//
//	if err != nil {
//		return nil,err
//	}
//	defer resp.Body.Close()
//
//	body, _ := ioutil.ReadAll(resp.Body)
//	logger.Debug(string(body))
//
//	var trasaResp models.TrasaResponse
//	err=json.Unmarshal([]byte(body), &trasaResp)
//
//	if err != nil {
//		return nil,err
//	}
//
//	return &trasaResp,nil
//
//
//}
//

func GetHomeDirAndUID() (string, int, int, error) {

	userc, err := user.Current()
	// fmt.Println(err)
	// fmt.Println(userc.HomeDir)

	if err != nil {
		return "", -1, -1, err
	}

	Context.OS_USERNAME = userc.Username
	uid, err := strconv.Atoi(userc.Uid)
	if err != nil {
		return userc.HomeDir, -1, -1, nil
	}
	gid, err := strconv.Atoi(userc.Gid)
	if err != nil {
		return userc.HomeDir, -1, -1, nil
	}

	return userc.HomeDir, uid, gid, nil

}

func GetHostConfig(setNewConf bool) {
	homeDir, _, _, err := GetHomeDirAndUID()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(0)
	}

	_, err = os.Stat(filepath.Join(homeDir, ".trasaconfig.toml"))
	if err != nil {
		//If file does not exist create one
		if errors.Is(err, os.ErrNotExist) {
			f, err := os.OpenFile(filepath.Join(homeDir, ".trasaconfig.toml"), os.O_CREATE, os.ModePerm)
			if err != nil {
				logger.Fatal(err)
			}
			f.Close()

		} else {
			logger.Fatal(err)
		}
	}

	viper.AddConfigPath(homeDir)
	viper.SetConfigName(".trasaconfig")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	trasaurl := viper.GetString("url")
	Context.TRASA_ORG_ID = viper.GetString("orgid")
	Context.KEYS_DIR_PATH = viper.GetString("keydirpath")
	Context.TRASA_ID = viper.GetString("trasaid")
	Context.SSH_CLIENT_NAME = viper.GetString("sshclient")
	if trasaurl == "" || setNewConf {
		result, err := promptTrasaURL.Run()
		if err != nil {
			logger.Fatal(err)
		}
		viper.Set("url", result)
		viper.WriteConfig()
		trasaurl = result
	}

	Context.TRASA_URL = trasaurl
	u, err := url.Parse(Context.TRASA_URL)
	if err != nil {
		logger.Fatal(err)
	}
	Context.TRASA_HOST = u.Host

	if Context.SSH_CLIENT_NAME == "" || setNewConf {
		//i, _, err := promptSSHClientName.Run()
		//if err != nil {
		//	logger.Fatal(err)
		//}
		//Context.SSH_CLIENT_NAME = sshclients[i].Name
		//

		Context.SSH_CLIENT_NAME = "openssh"

		viper.Set("sshclient", Context.SSH_CLIENT_NAME)
		viper.WriteConfig()
	}
	//fmt.Println("SSH Client: ", Context.SSH_CLIENT_NAME)
	//fmt.Println("TRASA URL: ", Context.TRASA_URL)

}

func GetLogLocation() string {
	homeDir, _, _, _ := GetHomeDirAndUID()
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/trasada.log"
	case "linux":
		return "/var/log/trasada.log"
	case "windows":
		return filepath.Join(homeDir, "trasada.log")
	default:
		return "trasada.log"
	}
}
