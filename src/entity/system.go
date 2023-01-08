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
	"sort"
	"time"
)

type SystemFSInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

type SystemNetworkInfo struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
	Rx   uint64 `json:"rx"`
	Tx   uint64 `json:"tx"`
}

type SystemCpuRaw struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

type SystemCPUInfo struct {
	User    float32 `json:"user"`
	Nice    float32 `json:"nice"`
	System  float32 `json:"system"`
	Idle    float32 `json:"idle"`
	Iowait  float32 `json:"iowait"`
	Irq     float32 `json:"irq"`
	SoftIrq float32 `json:"softIrq"`
	Steal   float32 `json:"steal"`
	Guest   float32 `json:"guest"`
}

type SystemStats struct {
	Uptime      time.Duration                 `json:"uptime"`
	Hostname    string                        `json:"hostname"`
	MemTotal    uint64                        `json:"memTotal"`
	MemFree     uint64                        `json:"memFree"`
	MemBuffers  uint64                        `json:"memBuffers"`
	MemCached   uint64                        `json:"memCached"`
	SwapTotal   uint64                        `json:"swapTotal"`
	SwapFree    uint64                        `json:"swapFree"`
	NetworkInfo map[string]*SystemNetworkInfo `json:"networkInfo"`
	CPU         map[string]*SystemCPUInfo     `json:"cpu"`
	TimeStamp   int64                         `json:"timeStamp"`
}

func (stats *SystemStats) ToString() string {
	used := stats.MemTotal - stats.MemFree - stats.MemBuffers - stats.MemCached
	var result = fmt.Sprintf(
		//%s%s%s%s up %s%s%s
		`
Memory:
    free    = %s%s%s
    used    = %s%s%s
    buffers = %s%s%s
    cached  = %s%s%s
    swap    = %s%s%s free of %s%s%s

`,
		escBrightWhite, fmtBytes(stats.MemFree), escReset,
		escBrightWhite, fmtBytes(used), escReset,
		escBrightWhite, fmtBytes(stats.MemBuffers), escReset,
		escBrightWhite, fmtBytes(stats.MemCached), escReset,
		escBrightWhite, fmtBytes(stats.SwapFree), escReset,
		escBrightWhite, fmtBytes(stats.SwapTotal), escReset,
	)

	if len(stats.CPU) > 0 {
		result += "CPU:\n"
		for k, v := range stats.CPU {
			result += fmt.Sprintf("%s :%s%.2f%s%% user, %s%.2f%s%% sys, %s%.2f%s%% nice, %s%.2f%s%% idle, %s%.2f%s%% iowait, %s%.2f%s%% hardirq, %s%.2f%s%% softirq, %s%.2f%s%% guest\n",
				k,
				escBrightWhite, v.User, escReset,
				escBrightWhite, v.System, escReset,
				escBrightWhite, v.Nice, escReset,
				escBrightWhite, v.Idle, escReset,
				escBrightWhite, v.Iowait, escReset,
				escBrightWhite, v.Irq, escReset,
				escBrightWhite, v.SoftIrq, escReset,
				escBrightWhite, v.Guest, escReset,
			)
		}
		result += "\n"
	}

	if len(stats.NetworkInfo) > 0 {
		result += "Network Interfaces:\n"
		keys := make([]string, 0, len(stats.NetworkInfo))
		for intf := range stats.NetworkInfo {
			keys = append(keys, intf)
		}
		sort.Strings(keys)
		for _, intf := range keys {
			info := stats.NetworkInfo[intf]
			result += fmt.Sprintf("    %s%s%s - %s%s%s",
				escBrightWhite, intf, escReset,
				escBrightWhite, info.IPv4, escReset,
			)
			if len(info.IPv6) > 0 {
				result += fmt.Sprintf(", %s%s%s\n",
					escBrightWhite, info.IPv6, escReset,
				)
			} else {
				result += "\n"
			}
			result += fmt.Sprintf("      rx = %s%s%s, tx = %s%s%s\n",
				escBrightWhite, fmtBytes(info.Rx), escReset,
				escBrightWhite, fmtBytes(info.Tx), escReset,
			)
			result += "\n"
		}
		result += "\n"
	}
	return result
}

func (stats *SystemStats) ToFormat() string {
	str, _ := json.MarshalIndent(stats, "", "\t")
	return string(str)
}

func (stats *SystemStats) ToJson() string {
	str, _ := json.Marshal(stats)
	return string(str)
}

const (
	escClear       = "\033[H\033[2J"
	escRed         = "\033[31m"
	escReset       = "\033[0m"
	escBrightWhite = "\033[37;1m"
)

func fmtBytes(val uint64) string {
	if val < 1024 {
		return fmt.Sprintf("%d bytes", val)
	} else if val < 1024*1024 {
		return fmt.Sprintf("%6.2f KiB", float64(val)/1024.0)
	} else if val < 1024*1024*1024 {
		return fmt.Sprintf("%6.2f MiB", float64(val)/1024.0/1024.0)
	} else {
		return fmt.Sprintf("%6.2f GiB", float64(val)/1024.0/1024.0/1024.0)
	}
}
