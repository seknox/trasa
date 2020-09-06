package main

import (
	"github.com/seknox/trasa/cli/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	cmd.Execute()
}
