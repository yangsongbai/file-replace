package tools

import (
	//logger "git.100tal.com/wangxiao_go_lib/xesLogger"
	"log"
	"os"
)

const (
	WINDOWS_OPERATOR_SYSTEM  = "windows"
)


func GetSystemDirSplit(osSystem string) string {
	if WINDOWS_OPERATOR_SYSTEM == osSystem {
		return "\\"
	} else {
		return  "/"
	}
}

func GetHostname() string {
	//hostname
	host, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return host
}
