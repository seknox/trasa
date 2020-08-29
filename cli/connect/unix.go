package connect

import (
	"fmt"
	"github.com/seknox/trasa/cli/config"
	"os"
	"os/exec"
)

func ConnectRemote(cmdType, user, port, certPath string, args ...string) {
	switch config.Context.SSH_CLIENT_NAME {
	case "openssh":
		connectUnix(cmdType, user, port, certPath, args...)
	case "putty":
		connectWindowsPutty(cmdType, user, port, certPath, args...)
	case "bitvise":
		connectWindowsBitvise(cmdType, user, port, certPath, args...)
	case "winscp":
		connectWindowsWinSCP(cmdType, user, port, certPath, args...)
	case "moba":
		connectWindowsMoba(cmdType, user, port, certPath, args...)
	}

}

func connectUnix(cmdType, user, port, certPath string, otherArgs ...string) {
	sshArgs := []string{
		"-o", fmt.Sprintf(`User=%s`, user),
		"-o", fmt.Sprintf(`Port=%s`, port),
		"-i", certPath,
	}

	sshArgs = append(sshArgs, otherArgs...)

	//scp takes hostname as <hostname>:<file path> so it is defered to scp itself
	if cmdType != "scp" {
		sshArgs = append(sshArgs, config.Context.TRASA_HOST)
	}

	//sshCmd := exec.Command(`ssh`, "-l", utils.Context.SSH_USERNAME, "-i", certPath, "localhost", "-p", "8022")
	c := exec.Command(cmdType, sshArgs...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

func connectWindowsPutty(cmdType, user, port, certPath string, otherArgs ...string) {

	args := []string{
		`-` + cmdType,
		fmt.Sprintf(`%s@%s`, user, config.Context.TRASA_HOST),
		`-P`, port,
		`-i`, certPath,
	}
	args = append(args, otherArgs...)

	puttyCmd := exec.Command(`putty.exe`,
		args...,
	)
	puttyCmd.Stderr = os.Stderr
	puttyCmd.Stdin = os.Stdin
	puttyCmd.Stdout = os.Stdout
	err := puttyCmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func connectWindowsBitvise(cmdType, user, port, certPath string, otherArgs ...string) {
	//C:\Program Files (x86)\Bitvise SSH Client\BvSsh.exe
	args := []string{
		`stermc`,
		fmt.Sprintf(`-host=%s`, config.Context.TRASA_HOST),
		fmt.Sprintf(`-port=%s`, port),
		fmt.Sprintf(`-user=%s`, user),
		fmt.Sprintf(`-keypairFile=%s`, certPath),
	}
	args = append(args, otherArgs...)

	puttyCmd := exec.Command(`C:\Program Files (x86)\Bitvise SSH Client\BvSsh.exe`,
		args...,
	)
	puttyCmd.Stderr = os.Stderr
	puttyCmd.Stdin = os.Stdin
	puttyCmd.Stdout = os.Stdout
	err := puttyCmd.Start()
	if err != nil {
		fmt.Println(err)
	}

}

func connectWindowsWinSCP(cmdType, user, port, certPath string, otherArgs ...string) {

}

func connectWindowsMoba(cmdType, user, port, certPath string, otherArgs ...string) {

}
