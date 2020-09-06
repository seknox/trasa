package config

import (
	"errors"
	"net/url"
	"strings"

	"github.com/manifoldco/promptui"
)

var validateURL = func(input string) error {
	u, err := url.Parse(input)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("Invalid URL")
	}
	return nil
}

var promptTrasaURL = promptui.Prompt{
	Label:    "Enter TRASA URL",
	Validate: validateURL,
}

type sshClient struct {
	Name string
	Path string
	Info string
}

var sshclients = []sshClient{
	sshClient{Name: "putty", Path: `putty.exe`, Info: ""},
	sshClient{Name: "moba", Path: `moba.exe`, Info: ""},
	sshClient{Name: "bitvise", Path: `bitvise.ext`, Info: ""},
	sshClient{Name: "winscp", Path: `winscp.exe`, Info: ""},
	sshClient{Name: "openssh", Path: `ssh`, Info: "Default ssh client of most unix systems"},
}

var promptSSHClientName = &promptui.Select{
	Label:        "Choose your default ssh client:",
	Items:        sshclients,
	IsVimMode:    false,
	HideHelp:     false,
	HideSelected: false,
	Templates: &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Info | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Info | red }})",
		Selected: "\U0001F3AF {{ .Name | red | cyan }}",
		Details: `
--------- Details ----------
{{ "Client Name:" | faint }}	{{ .Name }}
{{ "Path:" | faint }}	{{ .Path }}
{{ "Info:" | faint }}	{{ .Info }}`,
	},
	//Keys:              nil,
	Searcher: func(input string, index int) bool {
		c := sshclients[index].Name
		name := strings.Replace(strings.ToLower(c), " ", "", -1)
		ip := strings.Replace(strings.ToLower(c), " ", "", -1)

		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input) || strings.Contains(ip, input)
	},
	StartInSearchMode: true,
	//Pointer:           nil,
}
