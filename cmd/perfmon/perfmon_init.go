package perfmon

import (
	"github.com/spf13/cobra"
)

var perfmonRootCMD *cobra.Command

var interval int
var format bool
var serial string

func InitPerfmon(perfmonCMD *cobra.Command) {
	perfmonRootCMD = perfmonCMD
	initProcessPerfmon()
}
