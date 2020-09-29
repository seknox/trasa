package server_test

import (
	"fmt"
	"github.com/seknox/trasa/server/accessproxy/rdpproxy"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/providers/sidp"
	"github.com/seknox/trasa/server/api/providers/uidp"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/stats"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/initdb"
	"github.com/seknox/trasa/tests/server/testutils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestMain(m *testing.M) {
	state := setupTestEnv()
	logrus.SetLevel(logrus.TraceLevel)
	global.OxyLog = logrus.New()
	m.Run()
	tearDown(state)
}

func tearDown(state *global.State) {
	state.DB.Exec(`delete from users;`)
	state.DB.Exec(`delete from password_state;`)
	state.DB.Exec(`delete from org;`)
	state.DB.Exec(`delete from idp;`)
	state.DB.Exec(`delete from services;`)
	state.DB.Exec(`delete from adhoc_perms;`)
	state.DB.Exec(`delete from auth_logs;`)
	state.DB.Exec(`delete from backups;`)
	state.DB.Exec(`delete from org;`)
	state.DB.Exec(`delete from devices;`)
	state.DB.Exec(`delete from browsers;`)
	state.DB.Exec(`delete from cert_holder;`)
	state.DB.Exec(`delete from global_settings;`)
	state.DB.Exec(`delete from groups;`)
	state.DB.Exec(`delete from inapp_notifs;`)
	state.DB.Exec(`delete from key_holder;`)
	state.DB.Exec(`delete from keylog;`)
	state.DB.Exec(`delete from password_state;`)
	state.DB.Exec(`delete from policies;`)
	state.DB.Exec(`delete from security_rules;`)
	state.DB.Exec(`delete from user_accessmaps;`)
	state.DB.Exec(`delete from user_group_maps;`)
	state.DB.Exec(`delete from usergroup_accessmaps;`)

}

func setupTestEnv() *global.State {
	testConfig := global.Config{
		Database: struct {
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
		}{
			"postgres",
			"trasadb",
			"trasauser",
			"trasauser",
			"54321",
			"127.0.0.1",
			false,
			"", "", "",
		},
		Logging: struct {
			Level         string `toml:"level"`
			SendErrReport string `toml:"sendErrReport"`
		}{"TRACE", ""},
		Minio: struct {
			Status bool   `toml:"status"`
			Key    string `toml:"key"`
			Secret string `toml:"secret"`
			Server string `toml:"server"`
			Usessl bool   `toml:"usessl"`
		}{false, "", "", "", false},
		Platform: struct {
			Base string `toml:"base"`
		}{"private"},
		Redis: struct {
			Port       string   `toml:"port"`
			Server     []string `toml:"server"`
			Sslenabled bool     `toml:"sslenabled"`
			Usercert   string   `toml:"usercert"`
			Userkey    string   `toml:"userkey"`
			Cacert     string   `toml:"cacert"`
		}{"", []string{"127.0.0.1:16379"}, false, "", "", ""},

		Security: struct {
			InsecureSkipVerify bool `toml:"insecureSkipVerify"`
		}{true},
		Trasa: struct {
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
		}{false, "https://localhost", false, "localhost", "", "", "https://u2fproxy.trasa.io", "", "", testutils.MockOrgID},
		Proxy: struct {
			SSHListenAddr string `toml:"sshlistenAddr"`
			GuacdAddr     string `toml:"guacdAddr"`
			GuacdEnabled  bool   `toml:"guacdEnabled"`
		}{":8022", "127.0.0.1:4822", true},
	}

	state := global.InitDBSTOREWithConfig(testConfig)
	tearDown(state)
	err := insertMockData(state)
	if err != nil {
		panic(err)
	}

	rdpproxy.InitStore(state, accesscontrol.TrasaUAC)
	sshproxy.InitStore(state, accesscontrol.TrasaUAC)
	serviceauth.InitStore(state, accesscontrol.TrasaUAC)

	accesscontrol.InitStore(state, accesscontrol.TrasaUAC)

	accessmap.InitStore(state)

	auth.InitStore(state)

	tsxvault.InitStore(state)
	devices.InitStore(state)
	groups.InitStore(state)
	logs.InitStore(state)
	misc.InitStore(state)
	my.InitStore(state)
	notif.InitStore(state)
	orgs.InitStore(state)
	policies.InitStore(state)
	redis.InitStore(state)
	services.InitStore(state)
	system.InitStore(state)
	stats.InitStore(state)
	users.InitStore(state)
	ca.InitStore(state)
	sidp.InitStore(state)
	uidp.InitStore(state)
	vault.InitStore(state)

	initdb.InitDB()

	return state
}

func insertMockData(state *global.State) error {

	b, err := ioutil.ReadFile("mockdata.sql")
	if err != nil {
		return err
	}

	_, err2 := state.DB.Exec(string(b))
	if err2 != nil {
		fmt.Println(err2)
		return err2
	}

	return nil
}
