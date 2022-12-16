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
package perfmon

import (
	"bufio"
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
	"github.com/goinggo/mapstructure"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func getIoDataOnPid(client *adb.Device, pid string) (*entity.ProcessIO, error) {
	lines, err := client.OpenShell(fmt.Sprintf("/bin/cat /proc/%s/io", pid))
	if err != nil {
		return nil, fmt.Errorf("exec command erro : " + fmt.Sprintf("/bin/cat /proc/%s/io", pid))
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		log.Panic(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var ioMess = make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		value, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		ioMess[strings.TrimRight(fields[0], ":")] = value
	}
	var io = &entity.ProcessIO{}
	err = mapstructure.Decode(ioMess, io)
	//fmt.Println(lines)
	return io, nil
}

func getStatOnPid(client *adb.Device, pid string) (stat *entity.ProcessStat, err error) {
	lines, err := client.OpenShell(fmt.Sprintf("/bin/cat /proc/%s/stat", pid))
	if err != nil {
		return nil, fmt.Errorf("exec command erro : " + fmt.Sprintf("/bin/cat /proc/%s/stat", pid))
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		log.Panic(err)
	}
	return newProcessStat(string(data))
}

func GetPidOnAppName(client *adb.Device, appName string) (pid string, err error) {

	lines, err := client.OpenShell("/bin/ls /proc/")

	if err != nil {
		return "", fmt.Errorf("exec command erro : /bin/ls /proc/")
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		log.Panic(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		for _, part := range parts {
			if IsNum(part) {
				status, err := getStatusOnPid(client, part)
				if err != nil {
					continue
				}
				fmt.Println(status.Name, status.Pid)
				if status.Name == appName {
					return status.Pid, nil
				}
			}
		}

	}
	return "", fmt.Errorf("not find appname status")
}

func getStatusOnPid(client *adb.Device, pid string) (status *entity.ProcessStatus, err error) {
	lines, err1 := client.OpenShell(fmt.Sprintf("/bin/cat /proc/%s/status", pid))
	if err1 != nil {
		return status, fmt.Errorf("exec command erro : " + fmt.Sprintf("/bin/cat /proc/%s/status", pid))
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		log.Panic(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	status = &entity.ProcessStatus{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		var fieldName = strings.TrimRight(fields[0], ":")
		var value = strings.Join(fields[1:], " ")
		switch fieldName {
		case "Name":
			status.Name = value
		case "Umask":
			status.Umask = value
		case "State":
			status.State = value
		case "Tgid":
			status.Tgid = value
		case "Ngid":
			status.Ngid = value
		case "Pid":
			status.Pid = value
		case "PPid":
			status.PPid = value
		case "TracerPid":
			status.TracerPid = value
		case "Uid":
			status.Uid = value
		case "Gid":
			status.Gid = value
		case "FDSize":
			status.FDSize = value
		case "Groups":
			status.Groups = value
		case "VmPeak":
			status.VmPeak = value
		case "VmSize":
			status.VmSize = value
		case "VmLck":
			status.VmLck = value
		case "VmPin":
			status.VmPin = value
		case "VmHWM":
			status.VmHWM = value
		case "VmRSS":
			status.VmRSS = value
		case "RssAnon":
			status.RssAnon = value
		case "RssFile":
			status.RssFile = value
		case "RssShmem":
			status.RssShmem = value
		case "VmData":
			status.VmData = value
		case "VmStk":
			status.VmStk = value
		case "VmExe":
			status.VmExe = value
		case "VmLib":
			status.VmLib = value
		case "VmPTE":
			status.VmPTE = value
		case "VmSwap":
			status.VmSwap = value
		case "Threads":
			status.Threads = value
		case "SigQ":
			status.SigQ = value
		case "SigPnd":
			status.SigPnd = value
		case "ShdPnd":
			status.ShdPnd = value
		case "SigBlk":
			status.SigBlk = value
		case "SigIgn":
			status.SigIgn = value
		case "SigCgt":
			status.SigCgt = value
		case "CapInh":
			status.CapInh = value
		case "CapPrm":
			status.CapPrm = value
		case "CapEff":
			status.CapEff = value
		case "CapBnd":
			status.CapBnd = value
		case "CapAmb":
			status.CapAmb = value
		case "Cpus_allowed":
			status.CpusAllowed = value
		case "Cpus_allowed_list":
			status.CpusAllowedList = value
		case "voluntary_ctxt_switches":
			status.VoluntaryCtxtSwitches = value
		case "nonvoluntary_ctxt_switches":
			status.NonVoluntaryCtxtSwitches = value
		}
	}
	return status, err1
}

func newProcessStat(statStr string) (*entity.ProcessStat, error) {
	params := strings.Split(statStr, " ")
	var processStat = &entity.ProcessStat{}
	for i, value := range params {
		if i < 24 {
			switch i {
			case 0:
				processStat.Pid = value
			case 1:
				processStat.Comm = value
			case 2:
				processStat.State = value
			case 3:
				processStat.Ppid = value
			case 4:
				processStat.Pgrp = value
			case 5:
				processStat.Session = value
			case 6:
				processStat.Tty_nr = value
			case 7:
				processStat.Tpgid = value
			case 8:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Flags = num
				continue
			case 9:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Minflt = num
				continue
			case 10:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cminflt = num
			case 11:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Majflt = num
			case 12:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cmajflt = num
			case 13:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Utime = num
			case 14:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Stime = num
			case 15:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Cutime = num
			case 16:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Priority = num
			case 17:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Nice = num
			case 18:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Num_threads = num
			case 19:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Itrealvalue = num
			case 20:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Starttime = num
			case 21:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Vsize = num
			case 22:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Rss = num
			case 23:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Rss = num
			case 24:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Rsslim = num
			}
		}
	}
	return processStat, nil
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func getCpuUsage(client *adb.Device, pid string) {
	status, err := getStatOnPid(client, pid)
	if err != nil {
		fmt.Println(err)
		return
	}
	var nowTick = float64(status.Utime) + float64(status.Stime)
	if preTick == -1.0 {
		preTick = nowTick
		return
	}
	cpuUtilization = ((nowTick - preTick) / (HZ * sleepTime)) * 100
	preTick = nowTick
	time.Sleep(time.Duration(int(sleepTime) * int(time.Second)))
}

var preTick = -1.0
var sleepTime = 1.0 // # seconds
var HZ = 100.0      //# ticks/second
var cpuUtilization = 0.0

func GetProcessInfo(client *adb.Device, pid string, interval int64) (*entity.ProcessInfo, error) {
	sleepTime = float64(interval)

	stat, err := getStatOnPid(client, pid)
	if err != nil {
		return nil, err
	}
	status, err := getStatusOnPid(client, pid)
	if err != nil {
		return nil, err
	}
	//ioData, err := getIoDataOnPid(client, pid)
	//if err != nil {
	//	return nil, err
	//}

	var processInfo entity.ProcessInfo
	processInfo.PhyRSS = stat.Rss
	processInfo.VmSize = stat.Vsize
	if processInfo.Threads, err = strconv.Atoi(status.Threads); err != nil {
		return nil, err
	}

	getCpuUsage(client, pid)
	processInfo.CpuUtilization = cpuUtilization

	//processInfo.ReadBytes = ioData.ReadBytes
	//processInfo.WriteBytes = ioData.WriteBytes
	processInfo.Name = status.Name
	processInfo.Pid = status.Pid

	//processInfo.Rchar = ioData.Rchar
	//processInfo.Wchar = ioData.Wchar
	processInfo.TimeStamp = time.Now().UnixNano()
	return &processInfo, nil
}
