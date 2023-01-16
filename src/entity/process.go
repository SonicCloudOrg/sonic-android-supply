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
package entity

import (
	"encoding/json"
)

type ProcessIO struct {
	Rchar               int `json:"rchar"`
	Wchar               int `json:"wchar"`
	Syscr               int `json:"syscr"`
	Syscw               int `json:"syscw"`
	ReadBytes           int `json:"readBytes"`
	WriteBytes          int `json:"writeBytes"`
	CancelledWriteBytes int `json:"cancelledWriteBytes"`
}

type ProcessStat struct {
	Pid         string
	Comm        string
	State       string
	Ppid        string
	Pgrp        string
	Session     string
	Tty_nr      string
	Tpgid       string
	Flags       int
	Minflt      int
	Cminflt     int
	Majflt      int
	Cmajflt     int
	Utime       int
	Stime       int
	Cutime      int
	Cstime      int
	Priority    int
	Nice        int
	Num_threads int
	Itrealvalue int
	Starttime   int
	Vsize       int
	Rss         int
	Rsslim      int
}

type ProcessStatus struct {
	Name                     string `json:"name"`
	Umask                    string `json:"umask"`
	State                    string `json:"state"`
	Tgid                     string `json:"tgid"`
	Ngid                     string `json:"ngid"`
	Pid                      string `json:"pid"`
	PPid                     string `json:"pPid"`
	TracerPid                string `json:"tracerPid"`
	Uid                      string `json:"uid"`
	Gid                      string `json:"gid"`
	FDSize                   string `json:"fdSize"`
	Groups                   string `json:"groups"`
	VmPeak                   string `json:"vmPeak"`
	VmSize                   string `json:"vmSize"`
	VmLck                    string `json:"vmLck"`
	VmPin                    string `json:"vmPin"`
	VmHWM                    string `json:"vmHWM"`
	VmRSS                    string `json:"vmRSS"`
	RssAnon                  string `json:"rssAnon"`
	RssFile                  string `json:"rssFile"`
	RssShmem                 string `json:"rssShmem"`
	VmData                   string `json:"vmData"`
	VmStk                    string `json:"vmStk"`
	VmExe                    string `json:"vmExe"`
	VmLib                    string `json:"vmLib"`
	VmPTE                    string `json:"vmPTE"`
	VmSwap                   string `json:"vmSwap"`
	Threads                  string `json:"threads"`
	SigQ                     string `json:"sigQ"`
	SigPnd                   string `json:"sigPnd"`
	ShdPnd                   string `json:"shdPnd"`
	SigBlk                   string `json:"sigBlk"`
	SigIgn                   string `json:"sigIgn"`
	SigCgt                   string `json:"sigCgt"`
	CapInh                   string `json:"capInh"`
	CapPrm                   string `json:"capPrm"`
	CapEff                   string `json:"capEff"`
	CapBnd                   string `json:"capBnd"`
	CapAmb                   string `json:"capAmb"`
	CpusAllowed              string `json:"cpusAllowed"`
	CpusAllowedList          string `json:"cpusAllowedList"`
	VoluntaryCtxtSwitches    string `json:"voluntaryCtxtSwitches"`
	NonVoluntaryCtxtSwitches string `json:"nonVoluntaryCtxtSwitches"`
}

type ProcessInfo struct {
	Name           string   `json:"name"`
	Pid            string   `json:"pid"`
	CpuUtilization *float64 `json:"cpuUtilization,omitempty"`
	PhyRSS         *int     `json:"phyRSS,omitempty"`
	VmSize         *int     `json:"vmRSS,omitempty"`
	Threads        *int     `json:"threadCount,omitempty"`
	FPS            *int     `json:"fps,omitempty"`
	Error          []string `json:"error"`
}

func (i *ProcessInfo) ToJson() string {
	str, _ := json.Marshal(i)
	return string(str)
}

func (i *ProcessInfo) ToFormat() string {
	str, _ := json.MarshalIndent(i, "", "\t")
	return string(str)
}

func (i *ProcessInfo) ToString() string {
	return i.ToJson()
}

//func (i *ProcessInfo) ToString() string {
//
//	var result = fmt.Sprintf(
//		//%s%s%s%s up %s%s%s
//		`
//PID:%s%s%s Name:%s%s%s
//
//CPU:
//    %s%.2f%s%% cpuUtilizetion
//    %s%d%s     ThreadCount
//
//Memory:
//    physicalMemory    = %s%d%s
//    virtualMemory     = %s%d%s
//
//`,
//		escBrightWhite, i.Name, escReset,
//		escBrightWhite, i.Pid, escReset,
//		escBrightWhite, i.CpuUtilization, escReset,
//		escBrightWhite, i.Threads, escReset,
//		escBrightWhite, i.PhyRSS, escReset,
//		escBrightWhite, i.VmSize, escReset,
//	)
//	return result
//	//return fmt.Sprintf("name:%s pid:%s cpuUtilizetion:%f phyRss:%d vmRss:%d threadCount:%d readCharCount:%d writeCharCount:%d timeStamp:%d", i.Name, i.Pid, i.CpuUtilization, i.PhyRSS, i.VmSize, i.Threads, i.Rchar, i.Wchar, time.Now().Unix())
//}
