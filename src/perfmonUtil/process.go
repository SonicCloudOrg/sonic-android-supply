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
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
)

func getStatOnPid(client *adb.Device, pid string) (stat *entity.ProcessStat, err error) {
	lines, err := client.OpenShell(fmt.Sprintf("cat /proc/%s/stat", pid))
	if err != nil {
		return nil, fmt.Errorf("exec command erro : " + fmt.Sprintf("cat /proc/%s/stat", pid))
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		return
	}
	return newProcessStat(string(data))
}

func GetPidOnPackageName(client *adb.Device, appName string) (pid string, err error) {
	lines, err := client.OpenShell("dumpsys activity " + appName)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(fmt.Sprintf("ACTIVITY\\s+%s.*pid=[0-9]+", appName))
	regResult := reg.FindString(string(data))
	dataSplit := strings.Split(regResult, " ")
	if len(dataSplit) < 2 {
		return "", errors.New("unable to find the pid corresponding to app")
	}
	packageCurrentActivity = dataSplit[1]
	return strings.ReplaceAll(dataSplit[3], "pid=", ""), nil
}

func getMemTotalPSS(client *adb.Device, pid string) (result int, err error) {
	lines, err := client.OpenShell(fmt.Sprintf("dumpsys meminfo --local %s", pid))
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(lines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "TOTAL") {
			s := strings.Split(line, " ")
			flag := false
			for _, v := range s {
				if len(v) == 0 {
					continue
				}
				if v == "TOTAL" {
					flag = true
					continue
				}
				if flag {
					result, _ = strconv.Atoi(v)
					break
				}
			}
			break
		}
	}
	return
}

func getStatusOnPid(client *adb.Device, pid string) (status *entity.ProcessStatus, err error) {
	lines, err1 := client.OpenShell(fmt.Sprintf("cat /proc/%s/status", pid))
	if err1 != nil {
		return status, fmt.Errorf("exec command erro : " + fmt.Sprintf("cat /proc/%s/status", pid))
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	status = &entity.ProcessStatus{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
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
	status.TimeStamp = time.Now().UnixMilli()
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
				processStat.Majflt = num
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
				processStat.Cstime = num
			case 17:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Priority = num
			case 18:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Nice = num
			case 19:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Num_threads = num
			case 20:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Itrealvalue = num
			case 21:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Starttime = num
			case 22:
				num, err1 := strconv.Atoi(value)
				if err1 != nil {
					return nil, err1
				}
				processStat.Vsize = num
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
	processStat.TimeStamp = time.Now().UnixMilli()
	return processStat, nil
}

func getProcCpuUsage(stat *entity.ProcessStat) float64 {
	var nowTick = float64(stat.Utime) + float64(stat.Stime)
	if preTick == -1.0 {
		preTick = nowTick
		return 0.0
	}
	cpuUtilization := ((nowTick - preTick) / (HZ * IntervalTime)) * 100
	preTick = nowTick
	return cpuUtilization
}

var preTick = -1.0
var IntervalTime = 1.0 // # seconds
var HZ = 100.0         //# ticks/second
var packageCurrentActivity = ""
var PackageName = ""
var Pid = ""

func GetPIDAndPackageCurrentActivity(client *adb.Device, sign context.Context) {
	timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
	go func() {
		for {
			select {
			case <-sign.Done():
				return
			case <-timer:
				go func() {
					packageCurrentActivity = getPackageCurrentActivity(client, PackageName, Pid)
				}()
			}
		}
	}()
}

func GetNameOnPid(client *adb.Device, pid string) (string, error) {
	lines, err := client.OpenShell("ps -A")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(fmt.Sprintf(".*\\s+(%s)(\\s+\\d+){5}\\s\\S+\\s\\S+", pid))
	regResult := reg.FindString(string(data))
	reg = regexp.MustCompile("\\s+")
	regResult = reg.ReplaceAllString(regResult, " ")
	dataSplit := strings.Split(regResult, " ")
	if len(dataSplit) < 2 {
		return "", errors.New("unable to find the pid corresponding to app")
	}
	name := dataSplit[len(dataSplit)-1]
	return name, nil
}

func getPackageCurrentActivity(client *adb.Device, packageName string, pid string) string {
	lines, err := client.OpenShell("dumpsys activity " + packageName)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(lines)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(fmt.Sprintf("ACTIVITY\\s+%s.*%s", packageName, pid))
	regResult := reg.FindString(string(data))
	if regResult == "" {
		reg = regexp.MustCompile(fmt.Sprintf("ACTIVITY\\s+%s.*\\)", packageName))
		regResult = reg.FindString(string(data))
	}
	dataSplit := strings.Split(regResult, " ")
	if len(dataSplit) < 2 {
		return ""
	}
	if dataSplit[1] == "MANAGER" {
		return ""
	}
	return dataSplit[1]
}

func GetProcThreads(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.ProcThreads {
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						perfmonDataChan <- &entity.PerfmonData{
							Process: getThreads(client),
						}
					}()

				}
			}
		}()
	}
	return
}

func getThreads(client *adb.Device) *entity.ProcessInfo {
	processInfo := &entity.ProcessInfo{}
	status, err := getStatusOnPid(client, Pid)
	if err != nil {
		processInfo.Error = append(processInfo.Error, err.Error())
	}
	var threads int
	if len(status.Threads) > 0 {
		if threads, err = strconv.Atoi(status.Threads); err != nil {
			processInfo.Error = append(processInfo.Error, err.Error())
		}
	}
	processInfo.ThreadInfo = &entity.ProcTreadsInfo{
		Threads:   threads,
		TimeStamp: time.Now().UnixMilli(),
	}

	processInfo.Pid = Pid
	processInfo.Name = PackageName
	processInfo.Activity = packageCurrentActivity
	return processInfo
}

func GetProcFPS(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.ProcFPS {
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						perfmonDataChan <- &entity.PerfmonData{
							Process: getFPS(client),
						}
					}()

				}
			}
		}()
	}
	return
}

func getFPS(client *adb.Device) *entity.ProcessInfo {
	processInfo := &entity.ProcessInfo{}
	fpsInfo, err := getProcessFPSBySurfaceFlinger(client, PackageName)

	if fpsInfo.FPS <= 0 || err != nil {
		fpsInfo, err = getProcessFPSByGFXInfo(client, Pid)
	}
	if err != nil {
		processInfo.Error = append(processInfo.Error, err.Error())
	}
	processInfo.FPSInfo = &fpsInfo

	processInfo.Pid = Pid
	processInfo.Name = PackageName
	processInfo.Activity = packageCurrentActivity
	return processInfo
}

func GetProcCpu(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.ProcCPU {
		getProcCpu(client)
		time.Sleep(time.Duration(IntervalTime * float64(time.Second)))
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						perfmonDataChan <- &entity.PerfmonData{
							Process: getProcCpu(client),
						}
					}()

				}
			}
		}()
	}
	return
}

func getProcCpu(client *adb.Device) *entity.ProcessInfo {
	processInfo := &entity.ProcessInfo{}
	stat, err := getStatOnPid(client, Pid)
	if err != nil {
		processInfo.Error = append(processInfo.Error, err.Error())
	}

	processInfo.CPUInfo = &entity.ProcCpuInfo{
		CpuUtilization: getProcCpuUsage(stat),
		TimeStamp:      time.Now().UnixMilli(),
	}
	processInfo.Pid = Pid
	processInfo.Name = PackageName
	processInfo.Activity = packageCurrentActivity
	return processInfo
}

func GetProcMem(client *adb.Device, perfOptions entity.PerfOption, perfmonDataChan chan *entity.PerfmonData, sign context.Context) {
	if perfOptions.ProcMem {
		timer := time.Tick(time.Duration(int(IntervalTime * float64(time.Second))))
		go func() {
			for {
				select {
				case <-sign.Done():
					return
				case <-timer:
					go func() {
						perfmonDataChan <- &entity.PerfmonData{
							Process: getProcMem(client),
						}
					}()

				}
			}
		}()

	}
	return
}

func getProcMem(client *adb.Device) *entity.ProcessInfo {
	processInfo := &entity.ProcessInfo{}
	pss, _ := getMemTotalPSS(client, Pid)
	status, err := getStatusOnPid(client, Pid)
	if err != nil {
		processInfo.Error = append(processInfo.Error, err.Error())
	}
	vSize, _ := strconv.Atoi(strings.TrimSpace(strings.Split(status.VmSize," ")[0]))
	pRss, _ := strconv.Atoi(strings.TrimSpace(strings.Split(status.VmRSS," ")[0]))
	processInfo.MemInfo = &entity.ProcMemInfo{
		VmSize:    vSize,
		PhyRSS:    pRss,
		TotalPSS:  pss,
		TimeStamp: time.Now().UnixMilli(),
	}
	processInfo.Pid = Pid
	processInfo.Name = PackageName
	processInfo.Activity = packageCurrentActivity
	return processInfo
}

func getProcessFPSByGFXInfo(client *adb.Device, pid string) (result entity.ProcFPSInfo, err error) {
	lines, err := client.OpenShell(
		fmt.Sprintf("dumpsys gfxinfo %s | grep '.*visibility=0' -A129 | grep Draw -A128 | grep 'View hierarchy:' -B129", pid))
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(lines)
	frameCount := 0
	vsyncCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Draw") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		frameCount++
		s := strings.Split(line, "\t")
		if len(s) == 5 {
			render := RenderTime{}
			render.Draw, _ = strconv.ParseFloat(s[1], 64)
			render.Prepare, _ = strconv.ParseFloat(s[2], 64)
			render.Process, _ = strconv.ParseFloat(s[3], 64)
			render.Execute, _ = strconv.ParseFloat(s[4], 64)
			total := render.Draw + render.Prepare + render.Process + render.Execute

			if total > 16.67 {
				vsyncCount += (int)(math.Ceil(total/16.67) - 1)
			}
		}
	}
	result = entity.ProcFPSInfo{}
	if frameCount == 0 {
		result.FPS = 0
	} else {
		result.FPS = frameCount * 60 / (frameCount + vsyncCount)
	}
	result.TimeStamp = time.Now().UnixMilli()
	return
}

type RenderTime struct {
	Draw    float64
	Prepare float64
	Process float64
	Execute float64
}

func getProcessFPSBySurfaceFlinger(client *adb.Device, pkg string) (result entity.ProcFPSInfo, err error) {
	result = entity.ProcFPSInfo{}
	_, err = client.OpenShell("dumpsys SurfaceFlinger --latency-clear")
	lines, err := client.OpenShell(
		fmt.Sprintf("dumpsys SurfaceFlinger | grep %s", pkg))
	if err != nil {
		return
	}

	activity := ""

	scanner := bufio.NewScanner(lines)
	reg := regexp.MustCompile("\\[.*#0|\\(.*\\)")
	for scanner.Scan() {
		line := scanner.Text()

		activity = reg.FindString(line)

		if activity == "" {
			continue
		}
		break
	}
	if activity == "" {
		return result, errors.New(fmt.Sprintf("could not find app %s activity", pkg))
	}
	r := strings.NewReplacer("[", "", "(", "", ")", "")
	activity = r.Replace(activity)

	lines, err = client.OpenShell(
		fmt.Sprintf("dumpsys SurfaceFlinger --latency '%s'", activity))
	if err != nil {
		return
	}
	scanner = bufio.NewScanner(lines)
	var preFrame float64
	var t []float64
	for scanner.Scan() {
		line := scanner.Text()
		l := strings.Split(line, "\t")
		if len(l) < 3 {
			continue
		}
		if l[0][0] == '0' {
			continue
		}
		frame, _ := strconv.ParseFloat(l[1], 64)
		if frame == math.MaxInt64 {
			continue
		}
		frame /= 1e6
		if frame <= preFrame {
			continue
		}
		if preFrame == 0 {
			preFrame = frame
			continue
		}
		t = append(t, frame-preFrame)
		preFrame = frame
	}

	le := len(t)
	if le == 0 {
		return
	}
	result.FPS = (int)(float64(le) * 1000 / (sum(t, le)))
	result.TimeStamp = time.Now().UnixMilli()
	return
}

func sum(arr []float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	return sum(arr, n-1) + arr[n-1]
}
