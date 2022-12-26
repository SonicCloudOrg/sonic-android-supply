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
package perfmonUtil

import (
	"bufio"
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GetSystemStats(client *adb.Device) *entity.SystemStats {
	stats := &entity.SystemStats{}
	getAllStats(client, stats)
	return stats
}

func getAllStats(client *adb.Device, stats *entity.SystemStats) {
	getUptime(client, stats)
	getHostname(client, stats)
	getMemInfo(client, stats)
	getFSInfo(client, stats)
	getInterfaces(client, stats)
	getInterfaceInfo(client, stats)
	getCPU(client, stats)
	stats.TimeStamp = time.Now().UnixNano()
}

func getHostname(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/hostname -f")

	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		return
	}
	stats.Hostname = strings.TrimSpace(string(data))
	return
}

func getUptime(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/cat /proc/uptime")
	if err != nil {
		return
	}
	uptime, err := ioutil.ReadAll(lines)
	if err != nil {
		return
	}
	parts := strings.Fields(string(uptime))
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return
		}
		stats.Uptime = time.Duration(upsecs * 1e9)
	}

	return
}

func getMemInfo(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/cat /proc/meminfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				stats.MemTotal = val
			case "MemFree:":
				stats.MemFree = val
			case "Buffers:":
				stats.MemBuffers = val
			case "Cached:":
				stats.MemCached = val
			case "SwapTotal:":
				stats.SwapTotal = val
			case "SwapFree:":
				stats.SwapFree = val
			}
		}
	}
	return
}

func getFSInfo(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/df -B1")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(lines)
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		dev := n > 0 && strings.Index(parts[0], "/dev/") == 0
		if n == 1 && dev {
			flag = 1
		} else if (n == 5 && flag == 1) || (n == 6 && dev) {
			i := flag
			flag = 0
			used, err := strconv.ParseUint(parts[2-i], 10, 64)
			if err != nil {
				continue
			}
			free, err := strconv.ParseUint(parts[3-i], 10, 64)
			if err != nil {
				continue
			}
			stats.FSInfos = append(stats.FSInfos, &entity.SystemFSInfo{
				MountPoint: parts[5-i], Used: used, Free: free,
			})
		}
	}

	return
}

func getInterfaces(client *adb.Device, stats *entity.SystemStats) (err error) {
	var lines io.ReadCloser
	lines, err = client.OpenShell("/bin/ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = client.OpenShell("/sbin/ip -o addr")
		if err != nil {
			return
		}
	}

	if stats.NetworkInfo == nil {
		stats.NetworkInfo = make(map[string]*entity.SystemNetworkInfo)
	}

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := stats.NetworkInfo[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetworkInfo[intfname] = info
			} else {
				info := &entity.SystemNetworkInfo{}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetworkInfo[intfname] = info
			}
		}
	}

	return
}

func getInterfaceInfo(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/cat /proc/net/dev")
	if err != nil {
		return
	}

	if stats.NetworkInfo == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := stats.NetworkInfo[intf]; ok {
				rx, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					continue
				}
				tx, err := strconv.ParseUint(parts[9], 10, 64)
				if err != nil {
					continue
				}
				info.Rx = rx
				info.Tx = tx
				stats.NetworkInfo[intf] = info
			}
		}
	}
	return
}

func parseCPUFields(fields []string, stat *entity.SystemCpuRaw) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}

// the CPU stats that were fetched last time round
var preCPU entity.SystemCpuRaw

func getCPU(client *adb.Device, stats *entity.SystemStats) (err error) {
	lines, err := client.OpenShell("/bin/cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU entity.SystemCpuRaw
		total  float32
	)

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" { // changing here if want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
			break
		}
	}
	if preCPU.Total == 0 { // having no pre raw cpu data
		goto END
	}

	total = float32(nowCPU.Total - preCPU.Total)
	stats.CPU.User = float32(nowCPU.User-preCPU.User) / total * 100
	stats.CPU.Nice = float32(nowCPU.Nice-preCPU.Nice) / total * 100
	stats.CPU.System = float32(nowCPU.System-preCPU.System) / total * 100
	stats.CPU.Idle = float32(nowCPU.Idle-preCPU.Idle) / total * 100
	stats.CPU.Iowait = float32(nowCPU.Iowait-preCPU.Iowait) / total * 100
	stats.CPU.Irq = float32(nowCPU.Irq-preCPU.Irq) / total * 100
	stats.CPU.SoftIrq = float32(nowCPU.SoftIrq-preCPU.SoftIrq) / total * 100
	stats.CPU.Guest = float32(nowCPU.Guest-preCPU.Guest) / total * 100
END:
	preCPU = nowCPU
	return
}

func getFPS(client *adb.Device) (err error) {
	lines, err := client.OpenShell("dumpsys gfxinfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "jank") {
			fmt.Println("-============-")
		}
		fmt.Println(line)
	}
	return nil
}
