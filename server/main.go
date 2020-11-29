package main

import (
	"flag"
	"fmt"
	"github.com/seknox/trasa/server/global"
)

func main() {
	showversion := flag.Bool("v", false, "show version")
	global.LogToFile = flag.Bool("f", false, "Write to file")
	flag.Parse()
	if *showversion {
		fmt.Println("TRASA \nVersion: 1.1.4")
		return
	}
	StartServer()
}
