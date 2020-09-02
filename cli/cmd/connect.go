package cmd

import (
	"github.com/seknox/trasa/cli/config"
	"github.com/seknox/trasa/cli/connect"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to server through trasa",
	Long:  `Connect to server through trasa. <username> <hostname/IP/authapp name>`,
	Run:   runFuncGenerate("ssh"),
}

var sftpCmd = &cobra.Command{
	Use:   "sftp",
	Short: "Connect to server through trasa",
	Long:  `Connect to server through trasa. <username> <hostname/IP/authapp name>`,
	Run:   runFuncGenerate("sftp"),
}

var scpCmd = &cobra.Command{
	Use:   "scp",
	Short: "Connect to server through trasa",
	Long:  `Connect to server through trasa. <username> <hostname/IP/authapp name>`,
	Run:   runFuncGenerate("scp"),
}

var runFuncGenerate = func(cmdType string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		//Authenticate to trasa and download temporary certificate
		certPath := connect.Auth(config.Context.TRASA_ID, *config.Context.NEW_TRASA_ID)
		if certPath == "" {
			return
		}
		connect.ConnectRemote(cmdType, *config.Context.SSH_USERNAME, "8022", certPath, args...)

	}

}
