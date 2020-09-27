package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"unsafe"

	"github.com/sirupsen/logrus"
)

func main() {

	f, err := os.OpenFile(getLogLocation(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	defer func() {
		r := recover()
		if r != nil {
			logrus.Errorf("panic %v : stacktrace: %s", r, string(debug.Stack()))

		}
	}()

	logrus.SetOutput(f)

	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetReportCaller(true)

	// determine native byte order so that we can read message size correctly
	var one int16 = 1
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		nativeEndian = binary.BigEndian
	} else {
		nativeEndian = binary.LittleEndian
	}

	logrus.Debugf("trasaWrkstnAgent service started. Native byte order: %v.", nativeEndian)
	read()
	logrus.Debug("trasaWrkstnAgent service exited.")
}

func getLogLocation() string {
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/trasaWrkstnAgent.log"
	case "linux":
		return "/var/log/trasaWrkstnAgent.log"
	case "windows":
		usr, _ := user.Current()

		dir := fmt.Sprintf("%s\\AppData\\Local\\trasaWrkstnAgent", usr.HomeDir)
		createDirIfNotExist(dir)

		return filepath.Join(dir, "trasaWrkstnAgent.log")
	default:
		return "trasaWrkstnAgent.log"
	}
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)

	}
}
