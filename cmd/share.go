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
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share the connected adb device in the network",
	Long:  "Share the connected adb device in the network",
	Run: func(cmd *cobra.Command, args []string) {
		device := util.GetDevice(serial)

		adbd := adb.NewADBDaemon(device)
		fmt.Printf("Connect with port :%d\n", translatePort)
		err := adbd.ListenAndServe(":" + strconv.Itoa(translatePort))
		if err != nil {
			log.Panic(err)
		}
		return
	},
}

var translatePort int

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().IntVar(&translatePort, "translate-port", 6174, "translating proxy port")
	shareCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
}
