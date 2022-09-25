package cmd

import (
	"github.com/spf13/cobra"
	"go-android-supply/src/util"
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
