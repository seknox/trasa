package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/seknox/trasa/cli/config"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {

	isVerbose := rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	logToFile := rootCmd.PersistentFlags().BoolP("file log", "f", false, "log to file")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if *isVerbose {
			logger.SetLevel(logger.TraceLevel)
		} else {
			logger.SetLevel(logger.PanicLevel)
		}

		if *logToFile {

			f, err := os.OpenFile(getLogLocation(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				logger.Fatal(err)
			}
			defer f.Close()

			//This will output everything even fmt.println into file
			//os.Stdout=f
			//os.Stderr=f

			logger.SetOutput(f)
		}

		config.GetHostConfig(false)

	}

	var err error

	context.OS_TYPE = runtime.GOOS

	context.HOME_DIR, context.U_ID, context.G_ID, err = getHomeDirAndUID()
	if err != nil {
		fmt.Println(`Could not find home dir`)
	}

	config.Context.NEW_TRASA_ID = rootCmd.PersistentFlags().BoolP("new", "n", false, "Use new Trasa ID")
	config.Context.SSH_USERNAME = rootCmd.PersistentFlags().StringP("user", "u", config.Context.OS_USERNAME, "Upstream ssh username")
	config.Context.RAW_ARGS = rootCmd.PersistentFlags().StringP("args", "r", "", "Raw args to for")

	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(sftpCmd)
	rootCmd.AddCommand(scpCmd)
	rootCmd.AddCommand(configCmd)

}

var rootCmd = &cobra.Command{
	Use:   "trasa",
	Short: "TRASA protects your infrastructure with strong access control and advanced threat monitoring.",
	Long: `TRASA protects your infrastructure with strong access control and advanced threat monitoring.
Complete documentation is available at https://seknox.com/docs/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Trasa CLI. Version 0.0.1")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLogLocation() string {
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/fireser.log"
	case "linux":
		return "/var/log/fireser.log"
	case "windows":
		return filepath.Join(context.HOME_DIR, "fireser.log")
	default:
		return "fireser.log"
	}
}

var context struct {
	OS_TYPE         string
	SSH_CLIENT_NAME string
	HOME_DIR        string
	KEYS_DIR_PATH   string
	SSH_USERNAME    string
	OS_USERNAME     string
	U_ID            int
	G_ID            int
}

func getHomeDirAndUID() (string, int, int, error) {

	userc, err := user.Current()
	// fmt.Println(err)
	// fmt.Println(userc.HomeDir)

	if err != nil {
		return "", -1, -1, err
	}

	Context.OS_USERNAME = userc.Username
	uid, err := strconv.Atoi(userc.Uid)
	if err != nil {
		return userc.HomeDir, -1, -1, nil
	}
	gid, err := strconv.Atoi(userc.Gid)
	if err != nil {
		return userc.HomeDir, -1, -1, nil
	}

	return userc.HomeDir, uid, gid, nil

}
