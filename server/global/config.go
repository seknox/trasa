package global

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func parseConfig() (config Config) {
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
	Timezone struct {
		Location string `toml:"location"`
	} `toml:"timezone"`
	Security struct {
		InsecureSkipVerify bool `toml:"insecureSkipVerify"`
	} `toml:"security"`
	Trasa struct {
		ListenAddr  string `toml:"listenAddr"`
		Dashboard   string `toml:"dashboard"`
		Rootdomain  string `toml:"rootdomain"`
		CloudServer string `toml:"cloudServer"`
		Ssodomain   string `toml:"ssodomain"`
		Trasacore   string `toml:"trasacore"`
		Rootdir     string `toml:"rootdir"`
		OrgId       string `toml:"orgID"`
	} `toml:"trasa"`
	SSHProxy struct {
		ListenAddr string `toml:"listenAddr"`
	} `toml:"sshProxy"`
	Vault struct {
		Tsxvault bool   `toml:"tsxvault"`
		Port     string `toml:"port"`
		Server   string `toml:"server"`
		Token    string `toml:"token"`
	} `toml:"vault"`
	InternalHosts struct {
		Hosts string `toml:"hosts"`
	} `toml:"internalHosts"`
}
