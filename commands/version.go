package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number of emailctl",
	Long:  `Prints the version number of emailctl.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("emailctl v.0.1.0")
		fmt.Println("Postfix REST Server API V1 Client")
	},
}

func init() {
	emailctlCommand.AddCommand(versionCmd)
}
