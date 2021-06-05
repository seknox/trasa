package hygiene

/**
 * Copyright (C) 2020 Seknox Pte Ltd.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
)

var Context struct {
	OS_TYPE         string
	SSH_CLIENT_NAME string
	HOME_DIR        string
	KEYS_DIR_PATH   string
	SSH_USERNAME    string
	OS_USERNAME     string
	U_ID            int
	G_ID            int
}

func Pipe(c1, c2 *exec.Cmd) {
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
}

func GetRandomID(length int) (string, error) {
	val := make([]byte, length)
	_, err := rand.Read(val)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(val), nil

}

func GetHomeDirAndUID() (string, int, int, error) {

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

func GetLogLocation() string {
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/fireser.log"
	case "linux":
		return "/var/log/fireser.log"
	case "windows":
		return filepath.Join(Context.HOME_DIR, "fireser.log")
	default:
		return "fireser.log"
	}
}
