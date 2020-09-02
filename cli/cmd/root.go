package cmd

import (
	"fmt"
	"github.com/seknox/fireser/utils"
	"github.com/seknox/trasa/cli/config"
	"os"
	"runtime"

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

			f, err := os.OpenFile(utils.GetLogLocation(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	utils.Context.OS_TYPE = runtime.GOOS

	utils.Context.HOME_DIR, utils.Context.U_ID, utils.Context.G_ID, err = utils.GetHomeDirAndUID()
	if err != nil {
		fmt.Println(`Could not find home dir`)
	}

	config.Context.NEW_TRASA_ID = rootCmd.PersistentFlags().BoolP("new", "n", false, "Use new Trasa ID")
	config.Context.SSH_USERNAME = rootCmd.PersistentFlags().StringP("user", "u", config.Context.OS_USERNAME, "Upstream ssh username")

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
