/*
 *   sonic-android-supply  Supply of ADB.
 *   Copyright (C) 2022  SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package util

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"os"
	"os/exec"
	"regexp"
)

var (
	localADBHost = "127.0.0.1"
	localADBPort = 5037
)

func GetDevice(serial string) *adb.Device {
	client := adb.NewClient(fmt.Sprintf("%s:%d", localADBHost, localADBPort))
	if serial == "" {
		serialList, err := GetSerialList("")
		if err != nil {
			panic(err)
		}
		if len(serialList) == 0 {
			fmt.Println("failed to get serial list,not connect adb device")
			os.Exit(0)
		}
		serial = serialList[0]
	}
	device := client.DeviceWithSerial(serial)
	return device
}

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
