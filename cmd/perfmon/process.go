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
 *
 */
package perfmon

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmon"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var processPerfmonCmd = &cobra.Command{
	Use:   "process",
	Short: "get app or pid performance",
	Long:  "get app or pid performance",
	Run: func(cmd *cobra.Command, args []string) {
		if pid == "" && appName == "" {
			log.Println("pid or app-name is null")
			return
		}
		var err error
		device := util.GetDevice(serial)
		if pid == "" {
			// todo 优化
			pid, err = perfmon.GetPidOnAppName(device, appName)
			if err != nil {
				log.Panic(err)
			}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		timer := time.Tick(time.Duration(interval * int(time.Second)))
		done := false
		for !done {
			select {
			case <-sig:
				done = true
				fmt.Println()
			case <-timer:
				if processInfo, err := perfmon.GetProcessInfo(device, pid, 1); err != nil {
					log.Fatal(err)
				} else {
					if format {
						fmt.Println(processInfo.ToJson())
					} else {
						fmt.Println(processInfo.ToString())
					}
				}
			}
		}
		return
	},
}

var appName string
var pid string

func initProcessPerfmon() {
	perfmonRootCMD.AddCommand(processPerfmonCmd)
	processPerfmonCmd.Flags().StringVarP(&appName, "app-name", "n", "", "applicationName")
	processPerfmonCmd.Flags().StringVarP(&pid, "pid", "p", "", "process id")
	processPerfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
	processPerfmonCmd.Flags().IntVarP(&interval, "interval", "i", 1, "data refresh time")
	processPerfmonCmd.Flags().BoolVarP(&format, "format", "f", false, "formatted output")
}
