package cmd

import (
	"fmt"
	"github.com/codeskyblue/fa/adb"
	"github.com/spf13/cobra"
	"go-android-supply/src/util"
	"log"
	"strconv"
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Provides an USB device over TCP using a translating proxy.",
	Long:  "Provides an USB device over TCP using a translating proxy.",
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

func getSerial() {
	if serial == "" {
		serialList, err := util.GetSerialList("")
		if err != nil {
			log.Panic(err)
		}
		serial = serialList[0]
	}
}

var translatePort int

var serial string

var isTunnel bool

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().IntVar(&translatePort, "translate-port", 6174, "translating proxy port")
	translateCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
}
