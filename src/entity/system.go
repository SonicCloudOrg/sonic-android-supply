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
	IPv4 string
	IPv6 string
	Rx   uint64
	Tx   uint64
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
	User    float32
	Nice    float32
	System  float32
	Idle    float32
	Iowait  float32
	Irq     float32
	SoftIrq float32
	Steal   float32
	Guest   float32
}

type SystemStats struct {
	Uptime      time.Duration
	Hostname    string
	MemTotal    uint64
	MemFree     uint64
	MemBuffers  uint64
	MemCached   uint64
	SwapTotal   uint64
	SwapFree    uint64
	FSInfos     []*SystemFSInfo
	NetworkInfo map[string]*SystemNetworkInfo
	CPU         *SystemCPUInfo // or []SystemCPUInfo to get all the cpu-core's stats?
	TimeStamp   int64
}

func (stats *SystemStats) ToString() string {
	used := stats.MemTotal - stats.MemFree - stats.MemBuffers - stats.MemCached
	var result = fmt.Sprintf(
		//%s%s%s%s up %s%s%s
		`
CPU:
    %s%.2f%s%% user, %s%.2f%s%% sys, %s%.2f%s%% nice, %s%.2f%s%% idle, %s%.2f%s%% iowait, %s%.2f%s%% hardirq, %s%.2f%s%% softirq, %s%.2f%s%% guest

Memory:
    free    = %s%s%s
    used    = %s%s%s
    buffers = %s%s%s
    cached  = %s%s%s
    swap    = %s%s%s free of %s%s%s

`,
		escBrightWhite, stats.CPU.User, escReset,
		escBrightWhite, stats.CPU.System, escReset,
		escBrightWhite, stats.CPU.Nice, escReset,
		escBrightWhite, stats.CPU.Idle, escReset,
		escBrightWhite, stats.CPU.Iowait, escReset,
		escBrightWhite, stats.CPU.Irq, escReset,
		escBrightWhite, stats.CPU.SoftIrq, escReset,
		escBrightWhite, stats.CPU.Guest, escReset,
		escBrightWhite, fmtBytes(stats.MemFree), escReset,
		escBrightWhite, fmtBytes(used), escReset,
		escBrightWhite, fmtBytes(stats.MemBuffers), escReset,
		escBrightWhite, fmtBytes(stats.MemCached), escReset,
		escBrightWhite, fmtBytes(stats.SwapFree), escReset,
		escBrightWhite, fmtBytes(stats.SwapTotal), escReset,
	)
	if len(stats.FSInfos) > 0 {
		result += "Filesystems:\n"
		for _, fs := range stats.FSInfos {
			result += fmt.Sprintf("    %s%8s%s: %s%s%s free of %s%s%s\n",
				escBrightWhite, fs.MountPoint, escReset,
				escBrightWhite, fmtBytes(fs.Free), escReset,
				escBrightWhite, fmtBytes(fs.Used+fs.Free), escReset,
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
