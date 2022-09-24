package util

import (
	"net"
	"os/exec"
	"regexp"
)

func GetSerialList(adbPath string) (serialList []string, err error) {
	if adbPath == "" {
		adbPath = "adb"
	}
	output, err := exec.Command(adbPath, "devices", "-l").CombinedOutput()
	if err != nil {
		return
	}
	re := regexp.MustCompile(`(?m)^([^\s]+)\s+device\s+(.+)$`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	for _, m := range matches {
		serialList = append(serialList, m[1])
	}
	return
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}
