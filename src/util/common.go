/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
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
