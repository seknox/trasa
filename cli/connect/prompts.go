package connect

import (
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/seknox/trasa/server/models"
	"strconv"
	"strings"
)

var validateTOTP = func(input string) error {
	if len(input) == 0 {
		return nil
	}
	_, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errors.New("Invalid number")
	}
	if len(input) != 6 {
		return errors.New("Invalid length")
	}
	return nil
}

var promptTotp = promptui.Prompt{
	Label:    "TOTP",
	Validate: validateTOTP,
}

var nonEmpty = func(input string) error {
	if len(input) == 0 {
		return errors.New("Empty imput")
	}
	return nil
}

var promptUsername = promptui.Prompt{
	Label:    "Enter upstream username",
	Validate: nonEmpty,
}

var promptEmail = promptui.Prompt{
	Label:    "Enter TRASA email/username",
	Validate: nonEmpty,
}
var promptPassword = promptui.Prompt{
	Label:    "Enter TRASA password",
	Mask:     '*',
	Validate: nonEmpty,
}

func promptHostname(apps []models.MyServiceDetails) *promptui.Select {

	items := []string{}
	for _, i := range apps {
		if i.ServiceType == "ssh" {
			items = append(items, i.ServiceName+" : "+i.Hostname)
		}
	}

	return &promptui.Select{
		Label:        "Choose host to connect to:",
		Items:        apps,
		IsVimMode:    false,
		HideHelp:     false,
		HideSelected: false,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F449 {{ .AppName | cyan }} ({{ .Hostname | red }})",
			Inactive: "  {{ .AppName | cyan }} ({{ .Hostname | red }})",
			Selected: "\U0001F3AF {{ .AppName | red | cyan }}",
			Details: `
--------- Details ----------
{{ "AppName:" | faint }}	{{ .AppName }}
{{ "Hostname:" | faint }}	{{ .Hostname }}
{{ "Is2FAEnabled:" | faint }}	{{ .Is2FAEnabled }}`,
		},
		//Keys:              nil,
		Searcher: func(input string, index int) bool {
			app := apps[index]
			name := strings.Replace(strings.ToLower(app.ServiceName), " ", "", -1)
			ip := strings.Replace(strings.ToLower(app.Hostname), " ", "", -1)

			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input) || strings.Contains(ip, input)
		},
		StartInSearchMode: true,
		//Pointer:           nil,
	}

}
