package sshproxy

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

//This function takes  keyboardInteractive callback function to  get user to choose service/hostname
func chooseService(userEmail string, challengeUser ssh.KeyboardInteractiveChallenge) (*models.Service, error) {

	//take input(upstream server) from user and validate
	//loop until valid service is choosen
	for true {

		ans, err := challengeUser("user",
			"Choose Service",
			[]string{"\n\r_____________________________________________________________________________________\n\rEnter Service IP : \n\r"}, []bool{true})
		if len(ans) != 1 || err != nil {
			logrus.Debug("User canceled")
			return nil, fmt.Errorf("User canceled")
		}

		input := ans[0]
		isServiceName := false

		//check if the given input is service name or service hostname
		if strings.Contains(input, ":") {
			isServiceName = false
		} else {
			if isIP(input) {
				isServiceName = false
			} else {
				isServiceName = true
			}
		}

		var service *models.Service

		if isServiceName {
			service, err = services.Store.GetFromServiceName(input, global.GetConfig().Trasa.OrgId)
			if err != nil {
				challengeUser("", "Invalid service name", nil, nil)
				continue
			}
		} else {
			//If given input is hostname
			service, err = services.Store.GetFromHostname(input, "ssh", "", global.GetConfig().Trasa.OrgId)
			if errors.Is(err, sql.ErrNoRows) {
				service, err = accessmap.CreateDynamicService(input, "ssh", userEmail, global.GetConfig().Trasa.OrgId)
				if err != nil {
					logrus.Errorf("dynamic access: %v", err)
					challengeUser("", "Service not assigned. You do not have dynamic access", nil, nil)
					continue
				}

			} else if err != nil {
				logrus.Errorf("get service from hostname: %v", err)
				return nil, errors.WithMessage(err, "get service from hostname")
			}

		}

		h := service.Hostname
		if !strings.Contains(h, ":") {
			h = h + ":22"
		}

		//check if upstream server is down
		tempC, errPing := net.DialTimeout("tcp", h, time.Second*7)

		if errPing != nil {
			challengeUser("", "\n\nThe SSH server is down", nil, nil)
			continue
		}
		tempC.Close()
		challengeUser("", "\n\nService selected: "+service.Name+" \nHostname: "+service.Hostname+"\n\n", nil, nil)
		return service, nil
	}

	return nil, errors.New("Unexpected error")
}

func isIP(input string) bool {
	ip := net.ParseIP(input)
	return ip != nil
}
