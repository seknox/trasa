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
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

func Extract(reg string, from string) []string {
	re := regexp.MustCompile(reg)
	matches := re.FindStringSubmatch(from)
	return matches

}

func Execute(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		logrus.Trace(string(out))
		return "", err
	}
	str := string(NormalizeNewlines(out))
	return str, err
}

func ContainsIgnoreCase(s string, subStr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(subStr))
}
