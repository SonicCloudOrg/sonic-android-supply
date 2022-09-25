package cmd

import (
	"fmt"
	"github.com/codeskyblue/fa/adb"
	"github.com/spf13/cobra"
	"go-android-supply/src/perfmon"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var perfmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "get app performance",
	Long:  "get app performance",
	Run: func(cmd *cobra.Command, args []string) {
		if pid == "" && appName == "" {
			log.Println("pid or app-name is null")
			return
		}
		var err error
		getSerial()
		client := adb.NewClient(fmt.Sprintf("%s:%d", localADBHost, localADBPort))
		device := client.DeviceWithSerial(serial)
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
var interval int
var format bool

func init() {
	rootCmd.AddCommand(perfmonCmd)
	perfmonCmd.Flags().StringVarP(&appName, "app-name", "n", "", "applicationName")
	perfmonCmd.Flags().StringVarP(&pid, "pid", "p", "", "process id")
	perfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
	perfmonCmd.Flags().IntVarP(&interval, "interval", "i", 1, "data refresh time")
	perfmonCmd.Flags().BoolVarP(&format, "format", "f", false, "formatted output")
}
