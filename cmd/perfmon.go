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
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmonUtil"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
)

var perfmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Get device performance",
	Long:  "Get device performance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		device := util.GetDevice(serial)
		pidStr := ""
		if pid == -1 && packageName != "" {
			pidStr, err = perfmonUtil.GetPidOnPackageName(device, packageName)
			if err != nil {
				fmt.Println("no corresponding application PID found")
				os.Exit(0)
			}
		} else if pid != -1 && packageName == "" {
			pidStr = fmt.Sprintf("%d", pid)
			packageName, err = perfmonUtil.GetNameOnPid(device, pidStr)
			if err != nil {
				packageName = ""
			}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)

		if (pidStr == "" && packageName == "") &&
			!perfOptions.SystemCPU &&
			!perfOptions.SystemGPU &&
			!perfOptions.SystemNetWorking &&
			!perfOptions.SystemMem {
			sysAllParamsSet()
		}
		if (pidStr != "" || packageName != "") &&
			!perfOptions.ProcMem &&
			!perfOptions.ProcCPU &&
			!perfOptions.ProcThreads &&
			!perfOptions.ProcFPS {
			//sysAllParamsSet()
			perfOptions.ProcMem = true
			perfOptions.ProcCPU = true
			perfOptions.ProcThreads = true
			perfOptions.ProcFPS = true
		}
		perfmonUtil.PackageName = packageName
		perfmonUtil.Pid = pidStr

		exitCtx, exitChancel := context.WithCancel(context.Background())

		perfmonUtil.GetPIDAndPackageCurrentActivity(device, exitCtx)

		perfmonUtil.IntervalTime = float64(refreshTime) / 1000

		var perfDataChan = make(chan *entity.PerfmonData)
		perfmonUtil.GetSystemCPU(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetSystemMem(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetSystemNetwork(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcCpu(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcMem(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcFPS(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcThreads(device, perfOptions, perfDataChan, exitCtx)

		for {
			select {
			case <-sig:
				exitChancel()
				os.Exit(0)
			case perfData, ok := <-perfDataChan:
				if ok {
					fmt.Println(util.Format(perfData, isFormat, isJson))
				}
			}
		}
	},
}

var (
	perfOptions entity.PerfOption
	pid         int
	packageName string
	refreshTime int
)

func sysAllParamsSet() {
	perfOptions.SystemCPU = true
	perfOptions.SystemMem = true
	perfOptions.SystemGPU = true
	perfOptions.SystemNetWorking = true
}

func init() {
	rootCmd.AddCommand(perfmonCmd)
	perfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial (default first device)")
	perfmonCmd.Flags().IntVarP(&pid, "pid", "d", -1, "get PID data")
	perfmonCmd.Flags().StringVarP(&packageName, "package", "p", "", "app package name")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemCPU, "sys-cpu", false, "get system cpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemMem, "sys-mem", false, "get system memory data")
	//perfmonCmd.Flags().BoolVar(&sysDisk, "sys-disk", false, "get system disk data")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemNetWorking, "sys-network", false, "get system networking data")
	//perfmonCmd.Flags().BoolVar(&perfOptions.SystemGPU, "gpu", false, "get gpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcFPS, "proc-fps", false, "get fps data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcThreads, "proc-threads", false, "get process threads")
	//perfmonCmd.Flags().BoolVar(&, "proc-network", false, "get process network data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcCPU, "proc-cpu", false, "get process cpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcMem, "proc-mem", false, "get process mem data")
	perfmonCmd.Flags().IntVarP(&refreshTime, "refresh", "r", 1000, "data refresh time (millisecond)")
	perfmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	perfmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
}
