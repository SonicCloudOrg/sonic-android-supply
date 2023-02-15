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

type SystemNetworkInfo struct {
	IPv4      string `json:"ipv4"`
	IPv6      string `json:"ipv6"`
	Rx        uint64 `json:"rx"`
	Tx        uint64 `json:"tx"`
	TimeStamp int64  `json:"timeStamp"`
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
	User      float32 `json:"user"`
	Nice      float32 `json:"nice"`
	System    float32 `json:"system"`
	Idle      float32 `json:"idle"`
	Iowait    float32 `json:"iowait"`
	Irq       float32 `json:"irq"`
	SoftIrq   float32 `json:"softIrq"`
	Steal     float32 `json:"steal"`
	Guest     float32 `json:"guest"`
	Usage     float32 `json:"cpuUsage"`
	TimeStamp int64   `json:"timeStamp"`
}

type SystemMemInfo struct {
	MemTotal   uint64 `json:"memTotal"`
	MemFree    uint64 `json:"memFree"`
	MemBuffers uint64 `json:"memBuffers"`
	MemCached  uint64 `json:"memCached"`
	MemUsage   uint64 `json:"memUsage"`
	SwapTotal  uint64 `json:"swapTotal"`
	SwapFree   uint64 `json:"swapFree"`
	TimeStamp  int64  `json:"timeStamp"`
}

type SystemInfo struct {
	MemInfo     *SystemMemInfo                `json:"memInfo,omitempty"`
	NetworkInfo map[string]*SystemNetworkInfo `json:"networkInfo,omitempty"`
	CPU         map[string]*SystemCPUInfo     `json:"cpuInfo,omitempty"`
	Error       []string                      `json:"error"`
	//TimeStamp   int64                         `json:"timeStamp"`
}

func (stats *SystemInfo) ToString() string {
	return stats.ToJson()
}

func (stats *SystemInfo) ToFormat() string {
	str, _ := json.MarshalIndent(stats, "", "\t")
	return string(str)
}

func (stats *SystemInfo) ToJson() string {
	str, _ := json.Marshal(stats)
	return string(str)
}

type PerfOption struct {
	SystemCPU        bool
	SystemMem        bool
	SystemGPU        bool
	SystemNetWorking bool
	ProcCPU          bool
	ProcFPS          bool
	ProcMem          bool
	ProcThreads      bool
	RefreshTime      int
}
