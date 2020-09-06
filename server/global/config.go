package global

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ParseConfig uses viper to parse TRASA config file.
func ParseConfig() (config Config) {
	absPath, err := filepath.Abs("/etc/trasa/config")
	if err != nil {
		panic("config file not found in /etc/trasa/config")
	}
	//viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(absPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("config file invalid: " + err.Error())
	}

	if config.Minio.Key == "" {
		config.Minio.Key = "minioadmin"
	}
	if config.Minio.Secret == "" {
		config.Minio.Secret = "minioadmin"
	}

	//fmt.Println(config)
	//fmt.Println(viper.AllSettings())
	return
}

type Config struct {
	Backup struct {
		Backupdir string `toml:"backupdir"`
	} `toml:"backup"`
	Database struct {
		Dbtype     string `toml:"dbname"`
		Dbname     string `toml:"dbname"`
		Dbuser     string `toml:"dbuser"`
		Port       string `toml:"port"`
		Server     string `toml:"server"`
		Sslenabled bool   `toml:"sslenabled"`
		Usercert   string `toml:"usercert"`
		Userkey    string `toml:"userkey"`
		Cacert     string `toml:"cacert"`
	} `toml:"database"`
	Etcd struct {
		Server   string `toml:"server"`
		Usercert string `toml:"usercert"`
		Userkey  string `toml:"userkey"`
		Cacert   string `toml:"cacert"`
	} `toml:"etcd"`
	Logging struct {
		Level         string `toml:"level"`
		SendErrReport string `toml:"sendErrReport"`
	} `toml:"logging"`
	Minio struct {
		Status bool   `toml:"status"`
		Key    string `toml:"key"`
		Secret string `toml:"secret"`
		Server string `toml:"server"`
		Usessl bool   `toml:"usessl"`
	} `toml:"minio"`
	Platform struct {
		Base string `toml:"base"`
	} `toml:"platform"`
	Redis struct {
		Port       string   `toml:"port"`
		Server     []string `toml:"server"`
		Sslenabled bool     `toml:"sslenabled"`
		Usercert   string   `toml:"usercert"`
		Userkey    string   `toml:"userkey"`
		Cacert     string   `toml:"cacert"`
	} `toml:"redis"`

	Security struct {
		InsecureSkipVerify bool `toml:"insecureSkipVerify"`
	} `toml:"security"`
	Trasa struct {
		AutoCert    bool   `json:"autoCert"`
		ListenAddr  string `toml:"listenAddr"`
		Email       string `toml:"email"`
		Rootdomain  string `toml:"rootdomain"`
		CloudServer string `toml:"cloudServer"`
		Ssodomain   string `toml:"ssodomain"`
		Rootdir     string `toml:"rootdir"`
		OrgId       string `toml:"orgID"`
	} `toml:"trasa"`
	Proxy struct {
		SSHListenAddr string `toml:"sshlistenAddr"`
		GuacdAddr     string `toml:"guacdAddr"`
		GuacdEnabled  bool   `toml:"guacdEnabled"`
	} `toml:"proxy"`
	Vault struct {
		Tsxvault bool   `toml:"tsxvault"`
		Port     string `toml:"port"`
		Server   string `toml:"server"`
		Token    string `toml:"token"`
	} `toml:"vault"`
}

// UpdateTRASACPxyAddr updates TRASA cloud proxy server address.
func UpdateTRASACPxyAddr(serverAddr string) {
	absPath, err := filepath.Abs("/etc/trasa/config")
	if err != nil {
		logrus.Error(err)
	}
	//viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(absPath)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Error(err)
	}

	viper.Set("trasa.cloudServer", serverAddr)

	config.Trasa.CloudServer = serverAddr

	err = viper.WriteConfig()
	if err != nil {
		logrus.Error(err)
		return
	}

	// update global config val
	// config = ParseConfig()

	return
}
