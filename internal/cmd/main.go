package cmd

import (
	"github.com/spf13/cobra"
)

var appVersion string

func init() {
	addFlags()
	addCommands()
}

func Execute(version string) {
	appVersion = version
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())
}
