package global

import (
	"path/filepath"

	"github.com/seknox/trasa/server/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ParseConfig uses viper to parse TRASA config file.
func ParseConfig() (config Config) {
	absPath, err := filepath.Abs(filepath.Join(utils.GetETCDir(), "trasa", "config"))
	if err != nil {
		panic("config file not found in " + filepath.Join(utils.GetETCDir(), "trasa", "config"))
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
		Dbpass     string `toml:"dbpass"`
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
		ProxyDashboard bool   `toml:"proxyDashboard"`
		DashboardAddr  string `toml:"dashboardAddr"`
		AutoCert       bool   `toml:"autoCert"`
		ListenAddr     string `toml:"listenAddr"`
		Email          string `toml:"email"`
		Rootdomain     string `toml:"rootdomain"`
		CloudServer    string `toml:"cloudServer"`
		Ssodomain      string `toml:"ssodomain"`
		Rootdir        string `toml:"rootdir"`
		OrgId          string `toml:"orgID"`
	} `toml:"trasa"`
	Proxy struct {
		SSHListenAddr string `toml:"sshlistenAddr"`
		GuacdAddr     string `toml:"guacdAddr"`
		GuacdEnabled  bool   `toml:"guacdEnabled"`
	} `toml:"proxy"`
	Vault struct {
		Tsxvault bool   `toml:"tsxvault"`
		SaveMasterKey bool `toml:"saveMasterKey"`
		Key    string `toml:"key"`
	} `toml:"vault"`
}

// UpdateTRASACPxyAddr updates TRASA cloud proxy server address.
func UpdateTRASACPxyAddr(serverAddr string) {
	absPath, err := filepath.Abs(filepath.Join(utils.GetETCDir(), "trasa", "config"))
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


// SaveMasterKey saves master encryption key of tsxVault in config file. Only call this function if config value for "saveMasterKey" is true.
// This is meant to ease development flow only and should be false in production. 
func SaveMasterKey(key string) error {
	absPath, err := filepath.Abs(filepath.Join(utils.GetETCDir(), "trasa", "config"))
	if err != nil {
		return err
	}
	//viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(absPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	viper.Set("vault.key", key)

	config.Vault.Key = key

	err = viper.WriteConfig()
	if err != nil {
		return err

	}

	// update global config val
	// config = ParseConfig()

	return nil
}
