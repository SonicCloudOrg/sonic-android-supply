package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var localADBHost string

var localADBPort int

var rootCmd = &cobra.Command{
	Use:   "sab",
	Short: "Supply of Android Devices",
	Long:  ``,
}

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
