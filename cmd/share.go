package cmd

import (
	"fmt"
	"github.com/codeskyblue/fa/adb"
	"github.com/spf13/cobra"
	"go-android-supply/src/util"
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
