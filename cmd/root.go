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
	"github.com/spf13/cobra"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"log"
	"os"
)

var localADBHost string

var localADBPort int

var rootCmd = &cobra.Command{
	Use:   "sab",
	Short: "Supply of Android Devices",
	Long:  ``,
}

func getSerial() {
	if serial == "" {
		serialList, err := util.GetSerialList("")
		if err != nil {
			log.Panic(err)
		}
		serial = serialList[0]
	}
}

var serial string

// Execute error
func Execute() {

	localADBHost = "127.0.0.1"
	localADBPort = 5037

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	//err1 := doc.GenMarkdownTree(rootCmd, "doc")
	//if err1 != nil {
	//	log.Fatal(err1)
	//}
}
