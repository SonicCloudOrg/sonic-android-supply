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
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmonUtil"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var systemPerfmonCmd = &cobra.Command{
	Use:   "system",
	Short: "get system performance data",
	Long:  "get system performance data",
	Run: func(cmd *cobra.Command, args []string) {
		// todo
		device := util.GetDevice(serial)
		perfmonUtil.GetSystemStats(device)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		timer := time.Tick(time.Duration(interval * int(time.Second)))
		done := false
		for !done {
			select {
			case <-sig:
				done = true
				return
			case <-timer:
				status := perfmonUtil.GetSystemStats(device)
				data := util.ResultData(status)
				fmt.Println(util.Format(data, isFormat, isJson))
			}
		}
	},
}

func initSystemPerfmon() {
	perfmonRootCMD.AddCommand(systemPerfmonCmd)
	systemPerfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
	systemPerfmonCmd.Flags().IntVarP(&interval, "interval", "i", 1, "data refresh time")
	systemPerfmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	systemPerfmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
