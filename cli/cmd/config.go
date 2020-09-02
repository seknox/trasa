package cmd

import (
	"github.com/seknox/trasa/cli/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set new config",
	Long:  `Set new configs like TRASA URL, TRASA email and ssh client`,
	Run: func(cmd *cobra.Command, args []string) {
		config.GetHostConfig(true)
	},
}
