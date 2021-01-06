package accessproxytest

import (
	"net"
	"time"
)

func WaitUntilSSHServerStarts() {
	epoc := 1

	for epoc < 10 {
		c, err := net.Dial("tcp", "127.0.0.1:8022")
		if err == nil {
			c.Close()
			return
		}
		epoc++
		time.Sleep(time.Second)
	}

}
