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
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
)

func GetSystemCPU(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.SystemCPU {
		_ = getCPU(client, &entity.SystemInfo{})
		time.Sleep(time.Duration(IntervalTime * float64(time.Second)))
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						systemInfo := &entity.SystemInfo{}
						err := getCPU(client, systemInfo)
						if err != nil {
							systemInfo.Error = append(systemInfo.Error, err.Error())
						}
						if perfmonDataChan != nil {
							perfmonDataChan <- &entity.PerfmonData{
								System: systemInfo,
							}
						}
					}()

				}
			}
		}()
	}
	return
}

func GetSystemMem(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.SystemMem {
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						systemInfo := &entity.SystemInfo{}
						systemInfo.MemInfo = &entity.SystemMemInfo{}
						err := getMemInfo(client, systemInfo)
						if err != nil {
							systemInfo.Error = append(systemInfo.Error, err.Error())
						}
						if perfmonDataChan != nil {
							perfmonDataChan <- &entity.PerfmonData{
								System: systemInfo,
							}
						}
					}()

				}
			}
		}()

	}
	return
}

func GetSystemNetwork(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.SystemNetWorking {
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						systemInfo := &entity.SystemInfo{}
						err := getInterfaces(client, systemInfo)
						if err != nil {
							systemInfo.Error = append(systemInfo.Error, err.Error())
						}
						err = getInterfaceInfo(client, systemInfo)
						if err != nil {
							systemInfo.Error = append(systemInfo.Error, err.Error())
						}
						if perfmonDataChan != nil {
							perfmonDataChan <- &entity.PerfmonData{
								System: systemInfo,
							}
						}
					}()

				}
			}
		}()
	}
	return
}

func getMemInfo(client *adb.Device, stats *entity.SystemInfo) (err error) {
	lines, err := client.OpenShell("cat /proc/meminfo")
	stats.MemInfo.TimeStamp = time.Now().UnixMilli()
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
				stats.MemInfo.MemTotal = val
			case "MemFree:":
				stats.MemInfo.MemFree = val
			case "Buffers:":
				stats.MemInfo.MemBuffers = val
			case "Cached:":
				stats.MemInfo.MemCached = val
			case "SwapTotal:":
				stats.MemInfo.SwapTotal = val
			case "SwapFree:":
				stats.MemInfo.SwapFree = val
			}
		}
	}
	stats.MemInfo.MemUsage = stats.MemInfo.MemTotal - stats.MemInfo.MemFree - stats.MemInfo.MemBuffers - stats.MemInfo.MemCached
	return
}

func getInterfaces(client *adb.Device, stats *entity.SystemInfo) (err error) {
	var lines io.ReadCloser
	lines, err = client.OpenShell("ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = client.OpenShell("/bin/ip -o addr")
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

func getInterfaceInfo(client *adb.Device, stats *entity.SystemInfo) (err error) {
	lines, err := client.OpenShell("cat /proc/net/dev")
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
				info.TimeStamp = time.Now().UnixMilli()
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
var preCPUMap map[string]entity.SystemCpuRaw

func getCPU(client *adb.Device, stats *entity.SystemInfo) (err error) {
	lines, err := client.OpenShell("cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU entity.SystemCpuRaw
		total  float32
	)

	if preCPUMap == nil {
		preCPUMap = make(map[string]entity.SystemCpuRaw)
	}

	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") { // changing here if you want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
			preCPU = preCPUMap[fields[0]]
			if preCPU.Total == 0 { // having no pre raw cpu data
				preCPUMap[fields[0]] = nowCPU
				continue
			}

			total = float32(nowCPU.Total - preCPU.Total)
			if stats.CPU == nil {
				stats.CPU = make(map[string]*entity.SystemCPUInfo)
			}
			cpu := &entity.SystemCPUInfo{}
			cpu.User = float32(nowCPU.User-preCPU.User) / total * 100
			cpu.Nice = float32(nowCPU.Nice-preCPU.Nice) / total * 100
			cpu.System = float32(nowCPU.System-preCPU.System) / total * 100
			cpu.Idle = float32(nowCPU.Idle-preCPU.Idle) / total * 100
			cpu.Iowait = float32(nowCPU.Iowait-preCPU.Iowait) / total * 100
			cpu.Irq = float32(nowCPU.Irq-preCPU.Irq) / total * 100
			cpu.SoftIrq = float32(nowCPU.SoftIrq-preCPU.SoftIrq) / total * 100
			cpu.Guest = float32(nowCPU.Guest-preCPU.Guest) / total * 100
			var cpuNowTime = float32(nowCPU.User + nowCPU.Nice + nowCPU.System + nowCPU.Iowait + nowCPU.Irq + nowCPU.SoftIrq)
			var cpuPreTime = float32(preCPU.User + preCPU.Nice + preCPU.System + preCPU.Iowait + preCPU.Irq + preCPU.SoftIrq)

			cpu.Usage = (cpuNowTime - cpuPreTime) / ((cpuNowTime + float32(nowCPU.Idle)) - (cpuPreTime + float32(preCPU.Idle))) * 100
			cpu.TimeStamp = time.Now().UnixMilli()
			stats.CPU[fields[0]] = cpu
		}
	}
	return nil
}
