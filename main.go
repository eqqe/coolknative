package main

import (
	"os"

	"github.com/eskersoftware/coolknative/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdInstall := cmd.MakeInstall()
	cmdInfo := cmd.MakeInfo()


	var rootCmd = &cobra.Command{
		Use: "coolknative",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmdInstall)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdInfo)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
