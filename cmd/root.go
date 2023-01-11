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
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "sas",
	Short: "Supply of Android Devices",
	Long: `
   ▄▄▄▄      ▄▄▄▄    ▄▄▄   ▄▄   ▄▄▄▄▄▄      ▄▄▄▄
 ▄█▀▀▀▀█    ██▀▀██   ███   ██   ▀▀██▀▀    ██▀▀▀▀█
 ██▄       ██    ██  ██▀█  ██     ██     ██▀
  ▀████▄   ██    ██  ██ ██ ██     ██     ██
      ▀██  ██    ██  ██  █▄██     ██     ██▄
 █▄▄▄▄▄█▀   ██▄▄██   ██   ███   ▄▄██▄▄    ██▄▄▄▄█
  ▀▀▀▀▀      ▀▀▀▀    ▀▀   ▀▀▀   ▀▀▀▀▀▀      ▀▀▀▀

     Copyright (C) 2022 SonicCloudOrg AGPLv3
https://github.com/SonicCloudOrg/sonic-android-supply
`,
}

var (
	serial   string
	isFormat bool
	isJson   bool
)

// Execute error
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	//err1 := doc.GenMarkdownTree(rootCmd, "doc")
	//if err1 != nil {
	//	log.Fatal(err1)
	//}
}
