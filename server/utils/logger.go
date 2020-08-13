package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	logger "github.com/sirupsen/logrus"
)

func InitLogger() {
	logLevel := flag.String("l", "trace", "Set log level")
	logOutputToFile := flag.Bool("f", false, "Write to file")

	flag.Parse()
	level, _ := logger.ParseLevel(*logLevel)
	if *logOutputToFile {

		f, err := os.OpenFile("/var/log/trasa.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		logger.SetOutput(f)
	} else {
		logger.SetOutput(os.Stdout)
	}

	logger.SetReportCaller(true)

	logger.SetFormatter(&logger.TextFormatter{
		ForceColors:   false,
		DisableColors: false,
		//ForceQuote:                false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            true,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		//PadLevelText:              false,
		QuoteEmptyFields: false,
		FieldMap:         nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return filepath.Base(frame.Function), fmt.Sprintf(`%s:%d`, filepath.Base(frame.File), frame.Line)
		},
	})

	logger.SetLevel(level)
}

func MarshallStructStr(s interface{}) string {
	d, err := json.Marshal(s)
	if err != nil {
		logger.Debug(err)
	}
	return string(d)
}

//MarshallStructByte marshalls interface into bytes ignoring errors
func MarshallStructByte(s interface{}) []byte {
	d, err := json.Marshal(s)
	if err != nil {
		logger.Debug(err)
	}
	return d
}
