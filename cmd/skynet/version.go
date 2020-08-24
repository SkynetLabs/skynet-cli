package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of '%s'.", binName),
	Long:  fmt.Sprintf("Print the version number of '%s'.", binName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionString)
	},
}
