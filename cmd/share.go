/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package cmd

import (
	"fmt"
	"github.com/codeskyblue/fa/adb"
	"github.com/spf13/cobra"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"log"
	"strconv"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share the connected adb device in the network",
	Long:  "share the connected adb device in the network",
	Run: func(cmd *cobra.Command, args []string) {
		getSerial()
		client := adb.NewClient(fmt.Sprintf("%s:%d", localADBHost, localADBPort))
		device := client.DeviceWithSerial(serial)

		adbd := adb.NewADBDaemon(device)
		fmt.Printf("Connect with: adb connect %s:%d\n", util.GetLocalIP(), translatePort)
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
