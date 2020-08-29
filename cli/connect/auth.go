package connect

import (
	"fmt"
	"github.com/mholt/archiver"
	"github.com/seknox/trasa/cli/api"
	"github.com/seknox/trasa/cli/config"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//Auth will authenticate and download ssh certificate
func Auth(trasaID string, newTrasaID bool) (certPath string) {

	var err error
	if newTrasaID || trasaID == "" {
		trasaID, err = promptEmail.Run()
		if err != nil {
			fmt.Println("Invalid Email. ", err)
			return
		}
		viper.Set("trasaid", config.Context.TRASA_ID)
		viper.WriteConfig()
	}

	pass, err := promptPassword.Run()
	if err != nil {
		fmt.Println("Invalid Username")
		return
	}

	extCommCmd := exec.Command("/Users/bhrg3se/seknox/code/trasa/trasada/trasaExtNative/trasaExtNative", "get")
	out, err := extCommCmd.CombinedOutput()
	if err != nil {
		logger.Debug(err)
		fmt.Println("Could not get device hygiene trasaExtComm")
		return
	}
	fmt.Println(string(out))
	fmt.Println(trasaID, ":", pass)
	sshCertBytes, err := api.SendHygiene(trasaID, pass, string(out))
	if err != nil {
		fmt.Println("Could not update device hygiene")
		logger.Debug(err)
		return
	}

	homeDir, uid, gid, err := config.GetHomeDirAndUID()
	if err != nil {
		fmt.Println(`Could not find home dir`)
		logger.Fatal(err)
	}

	trasaKeysDir := filepath.Join(homeDir, ".ssh", "trasa_keys")

	os.MkdirAll(trasaKeysDir, 0700)

	//Remove old temp certs

	//Write new temp certificate to disk
	zipfile, err := ioutil.TempFile(trasaKeysDir, "*.zip")
	if err != nil {
		fmt.Println("Could not create temp file")
		logger.Fatal(err)
	}

	_, err = zipfile.Write(sshCertBytes)
	if err != nil {
		logger.Fatal(err)
	}
	defer zipfile.Close()
	z := archiver.NewZip()
	defer z.Close()
	os.Remove(filepath.Join(trasaKeysDir, "id_rsa"))
	os.Remove(filepath.Join(trasaKeysDir, "id_rsa.pub"))
	os.Remove(filepath.Join(trasaKeysDir, "id_rsa-cert.pub"))

	err = z.Unarchive(zipfile.Name(), filepath.Join(trasaKeysDir))
	if err != nil {
		fmt.Println(err)
		logger.Fatal(err)
	}

	os.Chown(filepath.Join(trasaKeysDir, "id_rsa"), uid, gid)
	os.Chmod(filepath.Join(trasaKeysDir, "id_rsa"), 0600)
	os.Remove(zipfile.Name())

	return filepath.Join(homeDir, ".ssh", "trasa_keys", "id_rsa")

}
