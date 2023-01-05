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
	"fmt"
)

type ProcessIO struct {
	Rchar               int `json:"rchar"`
	Wchar               int `json:"wchar"`
	Syscr               int `json:"syscr"`
	Syscw               int `json:"syscw"`
	ReadBytes           int `json:"read_bytes"`
	WriteBytes          int `json:"write_bytes"`
	CancelledWriteBytes int `json:"cancelled_write_bytes"`
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
	Name                     string `json:"Name"`
	Umask                    string `json:"Umask"`
	State                    string `json:"State"`
	Tgid                     string `json:"Tgid"`
	Ngid                     string `json:"Ngid"`
	Pid                      string `json:"Pid"`
	PPid                     string `json:"PPid"`
	TracerPid                string `json:"TracerPid"`
	Uid                      string `json:"Uid"`
	Gid                      string `json:"Gid"`
	FDSize                   string `json:"FDSize"`
	Groups                   string `json:"Groups"`
	VmPeak                   string `json:"VmPeak"`
	VmSize                   string `json:"VmSize"`
	VmLck                    string `json:"VmLck"`
	VmPin                    string `json:"VmPin"`
	VmHWM                    string `json:"VmHWM"`
	VmRSS                    string `json:"VmRSS"`
	RssAnon                  string `json:"RssAnon"`
	RssFile                  string `json:"RssFile"`
	RssShmem                 string `json:"RssShmem"`
	VmData                   string `json:"VmData"`
	VmStk                    string `json:"VmStk"`
	VmExe                    string `json:"VmExe"`
	VmLib                    string `json:"VmLib"`
	VmPTE                    string `json:"VmPTE"`
	VmSwap                   string `json:"VmSwap"`
	Threads                  string `json:"Threads"`
	SigQ                     string `json:"SigQ"`
	SigPnd                   string `json:"SigPnd"`
	ShdPnd                   string `json:"ShdPnd"`
	SigBlk                   string `json:"SigBlk"`
	SigIgn                   string `json:"SigIgn"`
	SigCgt                   string `json:"SigCgt"`
	CapInh                   string `json:"CapInh"`
	CapPrm                   string `json:"CapPrm"`
	CapEff                   string `json:"CapEff"`
	CapBnd                   string `json:"CapBnd"`
	CapAmb                   string `json:"CapAmb"`
	CpusAllowed              string `json:"Cpus_allowed"`
	CpusAllowedList          string `json:"Cpus_allowed_list"`
	VoluntaryCtxtSwitches    string `json:"voluntary_ctxt_switches"`
	NonVoluntaryCtxtSwitches string `json:"nonvoluntary_ctxt_switches"`
}

type ProcessInfo struct {
	Name           string  `json:"name"`
	Pid            string  `json:"pid"`
	CpuUtilization float64 `json:"cpuUtilization"`
	//ReadBytes      int     `json:"readBytes"`
	//WriteBytes     int     `json:"writeBytes"`
	PhyRSS    int   `json:"phyRSS"`
	VmSize    int   `json:"vmRSS"`
	Threads   int   `json:"threadCount"`
	Rchar     int   `json:"readCharCount"`
	Wchar     int   `json:"writeCharCount"`
	FPS       int   `json:"fps"`
	TimeStamp int64 `json:"timeStamp"`
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

	var result = fmt.Sprintf(
		//%s%s%s%s up %s%s%s
		`
PID:%s%s%s Name:%s%s%s

CPU:
    %s%.2f%s%% cpuUtilizetion
    %s%d%s     ThreadCount

Memory:
    physicalMemory    = %s%d%s
    virtualMemory     = %s%d%s

R/W char:
    Rchar = %s%d%s
	Wchar = %s%d%s

`,
		escBrightWhite, i.Name, escReset,
		escBrightWhite, i.Pid, escReset,
		escBrightWhite, i.CpuUtilization, escReset,
		escBrightWhite, i.Threads, escReset,
		escBrightWhite, i.PhyRSS, escReset,
		escBrightWhite, i.VmSize, escReset,
		escBrightWhite, i.Rchar, escReset,
		escBrightWhite, i.Wchar, escReset,
	)
	return result
	//return fmt.Sprintf("name:%s pid:%s cpuUtilizetion:%f phyRss:%d vmRss:%d threadCount:%d readCharCount:%d writeCharCount:%d timeStamp:%d", i.Name, i.Pid, i.CpuUtilization, i.PhyRSS, i.VmSize, i.Threads, i.Rchar, i.Wchar, time.Now().Unix())
}
